(() => {
  "use strict";

  const app = document.querySelector("#app");
  const searchInput = document.querySelector("#global-search");
  const connection = document.querySelector("#connection-status");
  const drawer = document.querySelector("#detail-drawer");
  const drawerContent = document.querySelector("#drawer-content");
  const drawerClose = document.querySelector("#drawer-close");

  const state = {
    campaigns: [],
    campaign: null,
    events: [],
    query: "",
    stream: null,
  };

  drawerClose.addEventListener("click", closeDrawer);
  searchInput.addEventListener("input", (event) => {
    state.query = event.target.value.trim().toLowerCase();
    renderRoute();
  });
  window.addEventListener("hashchange", loadRoute);
  window.addEventListener("keydown", (event) => {
    if (event.key === "/" && document.activeElement !== searchInput) {
      event.preventDefault();
      searchInput.focus();
    }
    if (event.key === "Escape") closeDrawer();
  });

  function escapeHTML(value) {
    return String(value ?? "")
      .replaceAll("&", "&amp;")
      .replaceAll("<", "&lt;")
      .replaceAll(">", "&gt;")
      .replaceAll('"', "&quot;")
      .replaceAll("'", "&#039;");
  }

  function shortID(value, length = 20) {
    const text = String(value || "—");
    if (text.length <= length) return text;
    const prefix = text.includes(":") ? `${text.split(":", 1)[0]}:` : "";
    const body = prefix ? text.slice(prefix.length) : text;
    return `${prefix}${body.slice(0, Math.max(6, length - prefix.length - 1))}…`;
  }

  function formatTime(value) {
    if (!value) return "—";
    const date = new Date(value);
    if (Number.isNaN(date.valueOf())) return String(value);
    return date.toISOString().replace("T", " ").replace(".000Z", "Z");
  }

  function statusClass(status) {
    if (status === "failed") return "failed";
    if (status === "running") return "running";
    if (status === "paused" || status === "stopping") return "paused";
    return "";
  }

  async function fetchJSON(url) {
    const response = await fetch(url, { headers: { Accept: "application/json" } });
    const body = await response.json().catch(() => ({}));
    if (!response.ok) {
      throw new Error(body?.error?.message || `${response.status} ${response.statusText}`);
    }
    return body;
  }

  function route() {
    const raw = location.hash.replace(/^#\/?/, "");
    const parts = raw.split("/").filter(Boolean).map(decodeURIComponent);
    if (parts[0] === "campaigns" && parts[1]) return { page: "campaign", id: parts[1] };
    return { page: "campaigns" };
  }

  async function loadRoute() {
    closeStream();
    closeDrawer();
    const current = route();
    app.innerHTML = '<section class="loading" aria-live="polite">READING CAMPAIGN JOURNAL…</section>';
    try {
      if (current.page === "campaign") {
        const [campaign, events] = await Promise.all([
          fetchJSON(`/api/v1/campaigns/${encodeURIComponent(current.id)}`),
          loadAllEvents(current.id),
        ]);
        state.campaign = campaign;
        state.events = events;
        renderCampaign();
        if (campaign.status === "running" || campaign.status === "paused" || campaign.status === "stopping") {
          openStream(current.id, campaign.version);
        } else {
          setConnection(campaign.integrity?.verified ? "VERIFIED" : "UNVERIFIED", campaign.integrity?.verified ? "live" : "error");
        }
      } else {
        const page = await fetchJSON("/api/v1/campaigns");
        state.campaigns = page.campaigns || [];
        state.campaign = null;
        state.events = [];
        setConnection("READ ONLY", "live");
        renderCampaignList();
      }
    } catch (error) {
      setConnection("QUERY ERROR", "error");
      renderError(error);
    }
  }

  async function loadAllEvents(id) {
    const events = [];
    let after = 0;
    for (;;) {
      const page = await fetchJSON(`/api/v1/campaigns/${encodeURIComponent(id)}/events?after=${after}&limit=500`);
      events.push(...(page.events || []));
      if (!page.has_more || page.through <= after) return events;
      after = page.through;
    }
  }

  function renderRoute() {
    if (state.campaign) renderCampaign();
    else renderCampaignList();
  }

  function renderCampaignList() {
    const query = state.query;
    const campaigns = state.campaigns.filter((campaign) => {
      if (!query) return true;
      return JSON.stringify(campaign).toLowerCase().includes(query);
    });
    app.innerHTML = `
      <div class="page-header">
        <div>
          <p class="eyebrow">SCIENTIFIC RECORD</p>
          <h1>Campaigns</h1>
          <p class="lede">CLI-authored experiments. Read-only evidence, replay, and provenance.</p>
        </div>
        <p class="section-count">${campaigns.length} / ${state.campaigns.length}</p>
      </div>
      <section class="section" aria-labelledby="campaign-list-title">
        <div class="section-heading">
          <h2 id="campaign-list-title">Local campaigns</h2>
          <span class="section-count">SEARCH WITH /</span>
        </div>
        <div class="campaign-list">
          ${campaigns.length ? campaigns.map(renderCampaignRow).join("") : '<p class="empty-state">No campaigns match this search.</p>'}
        </div>
      </section>`;
  }

  function renderCampaignRow(campaign) {
    return `
      <a class="campaign-row" href="#/campaigns/${encodeURIComponent(campaign.id)}">
        <span>
          <span class="campaign-name">${escapeHTML(campaign.system || "Campaign")}</span><br>
          <span class="entity-id">${escapeHTML(campaign.id)}</span>
        </span>
        <span class="mono">v${campaign.version}</span>
        <span class="status-text ${statusClass(campaign.status)}">${escapeHTML(campaign.status || "unknown").toUpperCase()}</span>
        <span>${campaign.overview?.completed || 0}/${campaign.overview?.episodes || 0}</span>
        <span>${campaign.integrity?.verified ? '<span class="status-text">VERIFIED</span>' : '<span class="error-text">CHECK</span>'}</span>
      </a>`;
  }

  function eventOf(kind) {
    return state.events.find((event) => event.kind === kind);
  }

  function eventsOf(kind) {
    return state.events.filter((event) => event.kind === kind);
  }

  function payloadOf(kind) {
    return eventOf(kind)?.payload_json || null;
  }

  function renderCampaign() {
    const campaign = state.campaign;
    if (!campaign) return;
    const candidatePayload = payloadOf("CandidateProposed");
    const candidate = candidatePayload?.candidate || null;
    const trial = payloadOf("TrialPlanned");
    const estimate = payloadOf("EstimateRecorded");
    const decision = payloadOf("DecisionRecorded");
    const filteredEvents = filterEvents(state.events, state.query);

    app.innerHTML = `
      <nav class="breadcrumbs" aria-label="Breadcrumb"><a href="#/campaigns">CAMPAIGNS</a> / ${escapeHTML(shortID(campaign.id, 26))}</nav>
      <div class="page-header">
        <div>
          <p class="eyebrow">${escapeHTML(campaign.system || "OPTIMIZATION CAMPAIGN")}</p>
          <h1>${escapeHTML(candidate?.hypothesis || "Campaign evidence")}</h1>
          <p class="lede entity-id">${escapeHTML(campaign.id)}</p>
        </div>
        <div>
          <p class="eyebrow">STATUS</p>
          <span class="status-text ${statusClass(campaign.status)}">${escapeHTML(campaign.status).toUpperCase()}</span>
        </div>
      </div>

      ${renderStageRail(campaign.event_kinds || {})}
      ${renderMetrics(campaign, estimate, decision)}
      ${renderLineage(candidatePayload, trial)}
      ${renderTrialMatrix(trial)}
      ${renderAnalysis(estimate, decision)}
      ${renderBudgets(campaign.budget)}
      ${renderEventKinds(campaign.event_kinds || {})}

      <section class="section" id="timeline" aria-labelledby="timeline-title">
        <div class="section-heading">
          <h2 id="timeline-title">Journal timeline</h2>
          <span class="section-count">${filteredEvents.length} / ${state.events.length} EVENTS</span>
        </div>
        <div class="event-tools">
          <span class="eyebrow">FILTERED BY GLOBAL SEARCH</span>
          ${state.query ? `<span class="mono">${escapeHTML(state.query)}</span>` : '<span class="mono">all events</span>'}
        </div>
        <div class="event-list" role="list">
          ${filteredEvents.map(renderEventRow).join("") || '<p>No matching events.</p>'}
        </div>
      </section>`;

    document.querySelectorAll("[data-event-index]").forEach((element) => {
      element.addEventListener("click", () => openEvent(Number(element.dataset.eventIndex)));
      element.addEventListener("keydown", (event) => {
        if (event.key === "Enter" || event.key === " ") openEvent(Number(element.dataset.eventIndex));
      });
    });
    document.querySelectorAll("[data-detail-kind]").forEach((element) => {
      element.addEventListener("click", () => openValue(element.dataset.detailKind, detailValue(element.dataset.detailKind, candidatePayload, trial, estimate, decision)));
    });
  }

  function renderStageRail(kinds) {
    const stages = [
      ["DESIGN", kinds.CampaignCreated && kinds.CandidateProposed],
      ["COMPILE", kinds.PlanCompiled && kinds.TrialPlanned],
      ["EXECUTE", kinds.EpisodeCompleted],
      ["MEASURE", kinds.ObservationRecorded],
      ["ANALYZE", kinds.EstimateRecorded],
      ["DECIDE", kinds.DecisionRecorded],
    ];
    return `<div class="stage-rail" aria-label="Campaign stages">${stages.map(([name, done]) => `<span class="stage ${done ? "done" : ""}">${name}</span>`).join("")}</div>`;
  }

  function renderMetrics(campaign, estimate, decision) {
    return `<section class="metric-grid" aria-label="Campaign metrics">
      ${metric("Journal", `v${campaign.version}`, `${campaign.overview?.events || 0} events`)}
      ${metric("Episodes", `${campaign.overview?.completed || 0}/${campaign.overview?.episodes || 0}`, `${campaign.overview?.failed_terminal || 0} failed`)}
      ${metric("Paired delta", estimate ? signed(estimate.value) : "—", estimate ? `${estimate.sample_size} pairs` : "no estimate")}
      ${metric("Decision", decision?.status?.toUpperCase() || "—", campaign.integrity?.verified ? "journal verified" : "integrity unknown")}
    </section>`;
  }

  function metric(label, value, note) {
    return `<div class="metric"><span class="metric-label">${escapeHTML(label)}</span><span class="metric-value">${escapeHTML(value)}</span><span class="metric-note">${escapeHTML(note)}</span></div>`;
  }

  function renderLineage(candidatePayload, trial) {
    const candidate = candidatePayload?.candidate;
    const patch = candidatePayload?.patch;
    if (!candidate) return "";
    return `<section class="section" aria-labelledby="lineage-title">
      <div class="section-heading"><h2 id="lineage-title">Configuration lineage</h2><span class="section-count">IMMUTABLE</span></div>
      <div class="lineage">
        ${lineageNode("BASELINE", candidate.parent, "baseline")}
        <span class="lineage-arrow" aria-hidden="true">→</span>
        ${lineageNode("PATCH", candidate.patch || patch?.id, "patch")}
        <span class="lineage-arrow" aria-hidden="true">→</span>
        ${lineageNode("CHALLENGER", candidate.child, "challenger")}
      </div>
      <p><span class="eyebrow">HYPOTHESIS</span>${escapeHTML(candidate.hypothesis)}</p>
      ${candidate.risks?.length ? `<p><span class="eyebrow warning-text">RECORDED RISK</span>${candidate.risks.map(escapeHTML).join("; ")}</p>` : ""}
      ${trial ? `<p class="entity-id">trial ${escapeHTML(trial.id || state.campaign.trial)}</p>` : ""}
    </section>`;
  }

  function lineageNode(label, value, kind) {
    return `<button class="lineage-node" type="button" data-detail-kind="${kind}"><span class="lineage-label">${label}</span><span class="lineage-value">${escapeHTML(value || "unavailable")}</span></button>`;
  }

  function renderTrialMatrix(trial) {
    if (!trial?.arms || !trial?.dataset?.cases) return "";
    const completions = new Map();
    for (const event of eventsOf("EpisodeCompleted")) {
      const item = event.payload_json;
      if (item?.case && item?.arm) completions.set(`${item.case}::${item.arm}::${item.repeat || 0}`, item);
    }
    const headers = trial.arms.map((arm) => `<th scope="col">${escapeHTML(arm.id)}</th>`).join("");
    const rows = trial.dataset.cases.map((testCase) => {
      const cells = trial.arms.map((arm) => {
        const completion = completions.get(`${testCase.id}::${arm.id}::0`);
        return `<td class="${completion ? "cell-complete" : "cell-queued"}">${completion ? `● COMPLETE<br><span class="entity-id">${escapeHTML(shortID(completion.episode, 18))}</span>` : "○ PLANNED"}</td>`;
      }).join("");
      return `<tr><th scope="row">${escapeHTML(testCase.id)}</th>${cells}</tr>`;
    }).join("");
    return `<section class="section" aria-labelledby="trial-title">
      <div class="section-heading"><h2 id="trial-title">Complete-block trial</h2><span class="section-count">${trial.dataset.cases.length} CASES × ${trial.arms.length} ARMS</span></div>
      <div class="table-wrap"><table><thead><tr><th>Case</th>${headers}</tr></thead><tbody>${rows}</tbody></table></div>
      <p class="entity-id">${escapeHTML(trial.id)} · ${escapeHTML(trial.protocol)} · ${trial.repeats} repeat</p>
    </section>`;
  }

  function renderAnalysis(estimate, decision) {
    if (!estimate && !decision) return "";
    const deltas = estimate?.deltas || [];
    const maximum = Math.max(1, ...deltas.map(Math.abs));
    const rows = deltas.map((delta, index) => `<div class="delta-row"><span class="mono">case-${String(index + 1).padStart(2, "0")}</span><progress class="delta-progress" max="${maximum}" value="${Math.abs(delta)}" aria-label="paired delta ${escapeHTML(signed(delta))}"></progress><span class="delta-value">${escapeHTML(signed(delta))}</span></div>`).join("");
    return `<section class="section" aria-labelledby="analysis-title">
      <div class="section-heading"><h2 id="analysis-title">Paired analysis</h2><button class="text-button" type="button" data-detail-kind="estimate">VIEW ESTIMATE JSON</button></div>
      <div class="delta-list">${rows}</div>
      ${decision ? `<p><span class="eyebrow ${decision.status === "eligible" ? "status-text" : "warning-text"}">${escapeHTML(decision.status)}</span>${escapeHTML(decision.reason || "")}</p><button class="text-button" type="button" data-detail-kind="decision">TRACE DECISION</button>` : ""}
    </section>`;
  }

  function renderBudgets(budget) {
    if (!budget?.resources) return "";
    const rows = budget.resources.map((resource) => {
      return `<div class="budget-row"><span class="mono">${escapeHTML(resource.resource)}</span><progress class="budget-progress" max="${Math.max(1, resource.limit)}" value="${resource.committed}" aria-label="${resource.committed} of ${resource.limit} committed"></progress><span class="budget-value">${resource.committed} / ${resource.limit}</span></div>`;
    }).join("");
    return `<section class="section" aria-labelledby="budget-title"><div class="section-heading"><h2 id="budget-title">Resource custody</h2><span class="section-count">${budget.violated ? '<span class="error-text">VIOLATED</span>' : '<span class="status-text">WITHIN LIMITS</span>'}</span></div><div class="budget-list">${rows}</div></section>`;
  }

  function renderEventKinds(kinds) {
    const entries = Object.entries(kinds).sort(([left], [right]) => left.localeCompare(right));
    return `<section class="section" aria-labelledby="facts-title"><div class="section-heading"><h2 id="facts-title">Fact inventory</h2><span class="section-count">CONTROL JOURNAL</span></div><div class="event-kind-grid">${entries.map(([kind, count]) => `<div class="event-kind-item"><span>${escapeHTML(kind)}</span><span>${count}</span></div>`).join("")}</div></section>`;
  }

  function filterEvents(events, query) {
    if (!query) return events;
    return events.filter((event) => JSON.stringify(event).toLowerCase().includes(query));
  }

  function renderEventRow(event) {
    const index = state.events.indexOf(event);
    return `<div class="event-row" role="listitem" tabindex="0" data-event-index="${index}"><span class="event-seq">${event.seq}</span><span class="event-kind">${escapeHTML(event.kind)}</span><span class="event-subject">${escapeHTML(shortID(event.subject || "—", 30))}</span><span class="event-schema">${escapeHTML(event.schema)}</span><time class="event-time">${escapeHTML(formatTime(event.recorded_at))}</time></div>`;
  }

  function detailValue(kind, candidatePayload, trial, estimate, decision) {
    if (kind === "baseline") return { id: candidatePayload?.candidate?.parent, role: "baseline snapshot" };
    if (kind === "patch") return candidatePayload?.patch || { id: candidatePayload?.candidate?.patch };
    if (kind === "challenger") return { id: candidatePayload?.candidate?.child, role: "challenger snapshot" };
    if (kind === "estimate") return estimate;
    if (kind === "decision") return decision;
    if (kind === "trial") return trial;
    return null;
  }

  function openEvent(index) {
    const event = state.events[index];
    if (!event) return;
    const details = [
      ["sequence", event.seq], ["event", event.id], ["kind", event.kind],
      ["schema", event.schema], ["subject", event.subject || "—"], ["actor", event.actor],
      ["occurred", formatTime(event.occurred_at)], ["recorded", formatTime(event.recorded_at)],
      ["digest", event.digest], ["payload", event.payload?.digest],
    ];
    openDrawer(`${event.kind} · ${event.seq}`, details, event.payload_available ? event.payload_json : { unavailable: event.payload_reason || "preview unavailable", ref: event.payload });
  }

  function openValue(title, value) {
    openDrawer(title.toUpperCase(), [], value);
  }

  function openDrawer(title, details, value) {
    drawerContent.innerHTML = `<p class="eyebrow">EVIDENCE DETAIL</p><h2>${escapeHTML(title)}</h2>${details.length ? `<dl class="detail-list">${details.map(([key, item]) => `<dt>${escapeHTML(key)}</dt><dd>${escapeHTML(item ?? "—")}</dd>`).join("")}</dl>` : ""}<h3>Payload</h3><pre id="drawer-json"></pre>`;
    document.querySelector("#drawer-json").textContent = JSON.stringify(value ?? null, null, 2);
    drawer.hidden = false;
    drawerClose.focus();
  }

  function closeDrawer() {
    drawer.hidden = true;
    drawerContent.innerHTML = "";
  }

  function signed(value) {
    const number = Number(value);
    if (!Number.isFinite(number)) return "—";
    return `${number >= 0 ? "+" : ""}${number.toFixed(6)}`;
  }

  function renderError(error) {
    const template = document.querySelector("#error-template").content.cloneNode(true);
    template.querySelector(".error-message").textContent = error.message || String(error);
    app.replaceChildren(template);
  }

  function openStream(id, after) {
    closeStream();
    setConnection("LIVE", "live");
    const stream = new EventSource(`/api/v1/campaigns/${encodeURIComponent(id)}/stream?after=${after}`);
    state.stream = stream;
    stream.addEventListener("campaign-event", async (message) => {
      const event = JSON.parse(message.data);
      if (!state.events.some((prior) => prior.seq === event.seq)) state.events.push(event);
      try {
        state.campaign = await fetchJSON(`/api/v1/campaigns/${encodeURIComponent(id)}`);
        renderCampaign();
      } catch (error) {
        setConnection("REFRESH ERROR", "error");
      }
    });
    stream.onerror = () => setConnection("RECONNECTING", "error");
    stream.onopen = () => setConnection("LIVE", "live");
  }

  function closeStream() {
    if (state.stream) state.stream.close();
    state.stream = null;
  }

  function setConnection(text, className = "") {
    connection.textContent = text;
    connection.className = `connection-status ${className}`.trim();
  }

  if (!location.hash) location.hash = "#/campaigns";
  loadRoute();
})();
