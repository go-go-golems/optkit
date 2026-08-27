# Numbergame Evidence Snapshot

Generated from Optkit revision `5db1815dabbf1150d59a0b7554e023df17cf860a`. The temporary local store was deleted after capture.

## Demo summary

```json
{
  "campaign": "campaign:b200ae2487f1eed38eabab8f24cc3cfd",
  "candidate": "candidate:7352fad2bfeac8caa485885fc767bb8cd031c6b960ce679a3658d228e64a44e9",
  "baseline": "snapshot:abaf1354710a5858ba594c40937e06bf535ae3d9754d2b43f646109521a0b2e1",
  "challenger": "snapshot:d0073b9e4ab569008367851937c9397363238c4dcb2b9d1d217f128981ea95ec",
  "patch": "patch:95d9160627d523946057f28f01fdcfa315e57a42accca582dfa8d38de7bcd491",
  "trial": "trial:3347f0684e4c8e061316284058e903cfc0e6618231294a474c0f688c5d3e6619",
  "estimate": {
    "baseline": "baseline",
    "treatment": "challenger",
    "value": 0.71041675,
    "sample_size": 4,
    "missing": 0,
    "deltas": [
      0.5,
      0.666667,
      0.8,
      0.875
    ]
  },
  "decision": {
    "policy": "numbergame.lexicographic/v1",
    "status": "eligible",
    "paired_accuracy_delta": 0.71041675,
    "all_interventions_exercised": true,
    "budget_violated": false,
    "reason": "Every intervention was exercised and paired target accuracy improved."
  },
  "overview": {
    "campaign": "campaign:b200ae2487f1eed38eabab8f24cc3cfd",
    "version": 81,
    "status": "completed",
    "events": 81,
    "episodes": 8,
    "queued": 0,
    "active": 0,
    "completed": 8,
    "failed_terminal": 0,
    "last_digest": "sha256:611db6e6a05c766ae913d496f42c92bd0eba7a846ac4a1e1d674013cee4eef5e"
  },
  "budget": {
    "campaign": "campaign:b200ae2487f1eed38eabab8f24cc3cfd",
    "resources": [
      {
        "resource": "episodes",
        "limit": 8,
        "reserved": 0,
        "committed": 8,
        "available": 0
      },
      {
        "resource": "numbergame.operations",
        "limit": 8,
        "reserved": 0,
        "committed": 8,
        "available": 0
      }
    ],
    "violated": false
  },
  "sample_trajectory": {
    "digest": "sha256:fc958c5445c1cdec48625d2a8e4dfe0ae2cce80dbe7072620fb45e0b37beae14",
    "media_type": "application/vnd.optkit.record+json",
    "schema": "schema:optkit.trajectory/v1",
    "size": 3027,
    "sensitivity": "internal"
  }
}
```

## Campaign inspection

```json
{
  "overview": {
    "campaign": "campaign:b200ae2487f1eed38eabab8f24cc3cfd",
    "version": 81,
    "status": "completed",
    "events": 81,
    "episodes": 8,
    "queued": 0,
    "active": 0,
    "completed": 8,
    "failed_terminal": 0,
    "last_digest": "sha256:611db6e6a05c766ae913d496f42c92bd0eba7a846ac4a1e1d674013cee4eef5e"
  },
  "budget": {
    "campaign": "campaign:b200ae2487f1eed38eabab8f24cc3cfd",
    "resources": [
      {
        "resource": "episodes",
        "limit": 8,
        "reserved": 0,
        "committed": 8,
        "available": 0
      },
      {
        "resource": "numbergame.operations",
        "limit": 8,
        "reserved": 0,
        "committed": 8,
        "available": 0
      }
    ],
    "violated": false
  },
  "event_kinds": {
    "BudgetReserved": 8,
    "CampaignCompleted": 1,
    "CampaignCreated": 1,
    "CampaignStarted": 1,
    "CandidateProposed": 1,
    "DecisionRecorded": 1,
    "EpisodeAttemptStarted": 8,
    "EpisodeCompleted": 8,
    "EpisodeLeaseGranted": 8,
    "EpisodeScheduled": 8,
    "EstimateRecorded": 1,
    "ObservationRecorded": 24,
    "PlanCompiled": 1,
    "SnapshotMaterialized": 1,
    "TrialPlanned": 1,
    "UsageCommitted": 8
  },
  "tail": [
    {
      "seq": 1,
      "kind": "CampaignCreated",
      "occurred_at": "2026-08-25T00:01:06.931364309Z",
      "payload_digest": "sha256:f0e596c7850384cb46260699348f73c67acbf1558f91481e8a4cde05d1d58c08"
    },
    {
      "seq": 2,
      "kind": "PlanCompiled",
      "occurred_at": "2026-08-25T00:01:06.932428983Z",
      "payload_digest": "sha256:9b54aa33e09c9c4d9841169a578713106b5e537108efbda0579ad631779190f0"
    },
    {
      "seq": 3,
      "kind": "CampaignStarted",
      "occurred_at": "2026-08-25T00:01:06.933472146Z",
      "payload_digest": "sha256:ac721d247a589d7629cadccbcc42e9d639b7fcc2d36e31eb06f54aef33c02e3b"
    },
    {
      "seq": 4,
      "kind": "CandidateProposed",
      "subject": "candidate:7352fad2bfeac8caa485885fc767bb8cd031c6b960ce679a3658d228e64a44e9",
      "occurred_at": "2026-08-25T00:01:06.93449157Z",
      "payload_digest": "sha256:be9694881bd5dbfd83e88b9bea0945df26f524d605e4dab124e8becab982b27e"
    },
    {
      "seq": 5,
      "kind": "SnapshotMaterialized",
      "subject": "snapshot:d0073b9e4ab569008367851937c9397363238c4dcb2b9d1d217f128981ea95ec",
      "occurred_at": "2026-08-25T00:01:06.935435268Z",
      "payload_digest": "sha256:dbba2aa7fe6e7639391597167e20d7f5ebb8c0881d0a6aa2d168fc3bff8d08b8"
    },
    {
      "seq": 6,
      "kind": "TrialPlanned",
      "subject": "trial:3347f0684e4c8e061316284058e903cfc0e6618231294a474c0f688c5d3e6619",
      "occurred_at": "2026-08-25T00:01:06.936101786Z",
      "payload_digest": "sha256:9b54aa33e09c9c4d9841169a578713106b5e537108efbda0579ad631779190f0"
    },
    {
      "seq": 7,
      "kind": "BudgetReserved",
      "subject": "work:2b975e226dc552703b388aaa27684f537f9eb2d3bb767e73c248491ccc2dd610",
      "occurred_at": "2026-08-25T00:01:06.938304226Z",
      "payload_digest": "sha256:163c9bd947eaf82df70f73fac6d724780ae81683e66d5ed71f6467f29a5444d3"
    },
    {
      "seq": 8,
      "kind": "EpisodeScheduled",
      "subject": "episode:12ab347f87c9009a82b6413b9c4078304387ffcf5b367f61f09b5a1ada910ba1",
      "occurred_at": "2026-08-25T00:01:06.938653747Z",
      "payload_digest": "sha256:ed35e21388ad475e79c1e2dce1088e82b27d3edd083989fe7dabb09dec63bd98"
    },
    {
      "seq": 9,
      "kind": "BudgetReserved",
      "subject": "work:9d11c53385838aa089e59161c9d2c43b68efa47e2bd4acddf84af6db9b4d0dc4",
      "occurred_at": "2026-08-25T00:01:06.94082002Z",
      "payload_digest": "sha256:28c79702e0ed8267c3f5630d33d1532cc251e596137c385e6f8886fd3ddf6d30"
    },
    {
      "seq": 10,
      "kind": "EpisodeScheduled",
      "subject": "episode:d61fe6b7eaf6ee06889a7d15b9084759c502692d265b43a1fd81ec4203c72a1f",
      "occurred_at": "2026-08-25T00:01:06.941068974Z",
      "payload_digest": "sha256:b231b0be07bcbf4636f2bb8dbcb5fb91f38bcaa76c30a170868f58d09fc4d04d"
    },
    {
      "seq": 11,
      "kind": "BudgetReserved",
      "subject": "work:c72c1d0d5e83480b41f7d31c4613bc630efe7c38eb4e41cb5af627a93fae25cc",
      "occurred_at": "2026-08-25T00:01:06.943050312Z",
      "payload_digest": "sha256:5c29b780cd597d1a244f5f95781503b719bc7e1bdae879079a236c1472ad60dd"
    },
    {
      "seq": 12,
      "kind": "EpisodeScheduled",
      "subject": "episode:f67a57c5a2e6dfe1d0d21c3543f59e981ad8caa505177f6f995d27dd2fc0ed2e",
      "occurred_at": "2026-08-25T00:01:06.943309379Z",
      "payload_digest": "sha256:8cbec007dbd9d5245954e4a4d2245adb06ab5afa7901172dd93c795c5a40a7c5"
    },
    {
      "seq": 13,
      "kind": "BudgetReserved",
      "subject": "work:ebaec2304a03a8f44079ea1f6020a1ecc94e35fdfb5215dd6fd7dfadb311fed5",
      "occurred_at": "2026-08-25T00:01:06.945462461Z",
      "payload_digest": "sha256:67ec4e9263d45f6314c845befa9a0ce27d6c05b2a302b5941c87176c8c7e8161"
    },
    {
      "seq": 14,
      "kind": "EpisodeScheduled",
      "subject": "episode:b4c54b1edb5ecec34d55170212c72e95f0a2f401224b1156f06581b69bd62ac1",
      "occurred_at": "2026-08-25T00:01:06.94578846Z",
      "payload_digest": "sha256:a1a7e3e77bcde8bdcf18adc2d02c87ec7129ea83fc1da3ffd24890a21c088062"
    },
    {
      "seq": 15,
      "kind": "BudgetReserved",
      "subject": "work:c2a260ab117fc8013bf6650143ff63be8063d262ea31040b3f181b0eafe1f897",
      "occurred_at": "2026-08-25T00:01:06.947752908Z",
      "payload_digest": "sha256:f9d4b017e5e85e9a00a05a9f2230069e6654bc205b5d3704e4c6ee12ca081d85"
    },
    {
      "seq": 16,
      "kind": "EpisodeScheduled",
      "subject": "episode:9e32c98093ea391c1c2e7ba89dda29f29ebad2cc1963175106c592b84d9948be",
      "occurred_at": "2026-08-25T00:01:06.948033961Z",
      "payload_digest": "sha256:3e1526b7a06cec7df04b21a1f6f0d28765cb620f462eaf5f48e27cfcfc71744f"
    },
    {
      "seq": 17,
      "kind": "BudgetReserved",
      "subject": "work:5c5b85f6a813204a583d51dca202027245ba65e20c80af913cb437a1ad113f3e",
      "occurred_at": "2026-08-25T00:01:06.950493547Z",
      "payload_digest": "sha256:132c1867d141c9d14b6ff3a6b45d2fedc1177d5d7c6c3cd6fa866d7d873b77d3"
    },
    {
      "seq": 18,
      "kind": "EpisodeScheduled",
      "subject": "episode:ed10febc1658b5170f8977165fcaa18254154326c01fe48b68b314ee4ac62fb5",
      "occurred_at": "2026-08-25T00:01:06.950965144Z",
      "payload_digest": "sha256:1700366dcf28f9cc359e98fbfb7b7b67f076ec06cfc738415ceb32eee43e846e"
    },
    {
      "seq": 19,
      "kind": "BudgetReserved",
      "subject": "work:adabf70aa0e7f84c4465f40eb9cab5ca05528879c1989147ffac336c374c78a2",
      "occurred_at": "2026-08-25T00:01:06.95301218Z",
      "payload_digest": "sha256:936c3582ba2fea573fac01fd57281a51ce223aa6af4f51c6ceaea3f8a439b627"
    },
    {
      "seq": 20,
      "kind": "EpisodeScheduled",
      "subject": "episode:a6eb91f286136e1eebb6b9c193c3bfd31377bd2d77c9ca2b14e7b93a028285f5",
      "occurred_at": "2026-08-25T00:01:06.953345618Z",
      "payload_digest": "sha256:c60049c3c0779d8db615f70546a516d5f5608cb3627141187a2c18193452ea89"
    },
    {
      "seq": 21,
      "kind": "BudgetReserved",
      "subject": "work:2c0685b50c91adab89e47180f469b660b712c00cbef4d3d734d1e448a064862b",
      "occurred_at": "2026-08-25T00:01:06.955296321Z",
      "payload_digest": "sha256:5ef0e077e2997fc756bef05594719b285a3388d2f763e2955fd0bb06cd124e1c"
    },
    {
      "seq": 22,
      "kind": "EpisodeScheduled",
      "subject": "episode:1b8d31e186897e58a3aa00fdf51a4893a2c8f867bd04e5af0d0d0bf1e81c5f7d",
      "occurred_at": "2026-08-25T00:01:06.955643825Z",
      "payload_digest": "sha256:0437c56cf09957f5a66b0cfa97fdd6aaf6569ba3cfa46e2f190a29c103ba5c25"
    },
    {
      "seq": 23,
      "kind": "EpisodeLeaseGranted",
      "subject": "episode:12ab347f87c9009a82b6413b9c4078304387ffcf5b367f61f09b5a1ada910ba1",
      "occurred_at": "2026-08-25T00:01:06.957217941Z",
      "payload_digest": "sha256:e963898e0ef6fc2df79a52f53481d7069b27bcbfbc791080d73a040996fed098"
    },
    {
      "seq": 24,
      "kind": "EpisodeAttemptStarted",
      "subject": "episode:12ab347f87c9009a82b6413b9c4078304387ffcf5b367f61f09b5a1ada910ba1",
      "occurred_at": "2026-08-25T00:01:06.958487041Z",
      "payload_digest": "sha256:4f997ec5bf69fdf6b76c5426c0f7b81e405a4cca0bae4184273367b607c09c93"
    },
    {
      "seq": 25,
      "kind": "UsageCommitted",
      "subject": "work:2b975e226dc552703b388aaa27684f537f9eb2d3bb767e73c248491ccc2dd610",
      "occurred_at": "2026-08-25T00:01:06.969932736Z",
      "payload_digest": "sha256:25869683008c87dc979aa4da68dc742907bb0030810db153c7fd8ebfe94050a7"
    },
    {
      "seq": 26,
      "kind": "EpisodeCompleted",
      "subject": "episode:12ab347f87c9009a82b6413b9c4078304387ffcf5b367f61f09b5a1ada910ba1",
      "occurred_at": "2026-08-25T00:01:06.970341739Z",
      "payload_digest": "sha256:7f87f0f6c94b36635c7ea100a520f55fffe1a1f8958a0c0d92f798bdd16f294b"
    },
    {
      "seq": 27,
      "kind": "ObservationRecorded",
      "subject": "observation:e79743dab8fd4b1fe6ff16a55640fb157df2ba49d375bb538b9db44332b434d8",
      "occurred_at": "2026-08-25T00:01:06.970727922Z",
      "payload_digest": "sha256:298e740d6e305f02d04a5bc55d055e09e04683b504f89b5128f25bb59791da25"
    },
    {
      "seq": 28,
      "kind": "ObservationRecorded",
      "subject": "observation:0f0d5ba1f174bc47300b200e606fc7befea319e1e630d50dce644ac5259d0120",
      "occurred_at": "2026-08-25T00:01:06.971160973Z",
      "payload_digest": "sha256:20b7bbb8dff59a84f4000b8a37a79e99b3e608baf3bdcf2b9a69d528dffeb5a2"
    },
    {
      "seq": 29,
      "kind": "ObservationRecorded",
      "subject": "observation:1159e0ae4eaee0a4eb1bbac61017e8f59b2c2ead399ab6771c2e14456b3f1f35",
      "occurred_at": "2026-08-25T00:01:06.971547513Z",
      "payload_digest": "sha256:5832023b0158f16b9aad5614b93da727f2b9f527eb1ad7a0cc39eed7b4ea197b"
    },
    {
      "seq": 30,
      "kind": "EpisodeLeaseGranted",
      "subject": "episode:d61fe6b7eaf6ee06889a7d15b9084759c502692d265b43a1fd81ec4203c72a1f",
      "occurred_at": "2026-08-25T00:01:06.972820074Z",
      "payload_digest": "sha256:327cf9b1cce2bc2671f1502c5161c52afca3ca418416b2df7fb5db79abf96aa5"
    },
    {
      "seq": 31,
      "kind": "EpisodeAttemptStarted",
      "subject": "episode:d61fe6b7eaf6ee06889a7d15b9084759c502692d265b43a1fd81ec4203c72a1f",
      "occurred_at": "2026-08-25T00:01:06.973926225Z",
      "payload_digest": "sha256:8c1a3408bca1d04ba33ac05e97705e2ac0f5c537daa313ff778d1c144ff7a429"
    },
    {
      "seq": 32,
      "kind": "UsageCommitted",
      "subject": "work:9d11c53385838aa089e59161c9d2c43b68efa47e2bd4acddf84af6db9b4d0dc4",
      "occurred_at": "2026-08-25T00:01:06.984035076Z",
      "payload_digest": "sha256:d16264753445b75f211e81d4bb52ce3876c21d342f77ae86cd392676a6226e45"
    },
    {
      "seq": 33,
      "kind": "EpisodeCompleted",
      "subject": "episode:d61fe6b7eaf6ee06889a7d15b9084759c502692d265b43a1fd81ec4203c72a1f",
      "occurred_at": "2026-08-25T00:01:06.984689265Z",
      "payload_digest": "sha256:b467c0f6f6f18413f5cf2ab9121407360b28ad4a4de5decf3cfdc1389d7d2eaf"
    },
    {
      "seq": 34,
      "kind": "ObservationRecorded",
      "subject": "observation:7b662a17eea4b3b425f59dcd83198a702717d6caba5c1a021ea72c10443a2518",
      "occurred_at": "2026-08-25T00:01:06.985148719Z",
      "payload_digest": "sha256:a57d02789c3c507b90cdfd437a4dfcb42e897d6b663970642b6f8cb74c798e30"
    },
    {
      "seq": 35,
      "kind": "ObservationRecorded",
      "subject": "observation:57167f672bb80928119d8744e6a287a51a3035d435a57cb22b52263b04178176",
      "occurred_at": "2026-08-25T00:01:06.985615252Z",
      "payload_digest": "sha256:b10ab814c31d9c2781ab457f0c33ece8dab7933fb8c849e3450114b8add89318"
    },
    {
      "seq": 36,
      "kind": "ObservationRecorded",
      "subject": "observation:21bbae38a80fec2585db39dbb75a63652a48948edafc4cf683f30279cc28eefb",
      "occurred_at": "2026-08-25T00:01:06.986069178Z",
      "payload_digest": "sha256:ef5a10532aa4a387ed45aa5bc2aa963a9cc90a67e2370f68dec4b8025e70e4e7"
    },
    {
      "seq": 37,
      "kind": "EpisodeLeaseGranted",
      "subject": "episode:f67a57c5a2e6dfe1d0d21c3543f59e981ad8caa505177f6f995d27dd2fc0ed2e",
      "occurred_at": "2026-08-25T00:01:06.98755144Z",
      "payload_digest": "sha256:9f7895b5de135af6f7563b4a16bdc63279a076efbd241f631dca23abd6d5fe3a"
    },
    {
      "seq": 38,
      "kind": "EpisodeAttemptStarted",
      "subject": "episode:f67a57c5a2e6dfe1d0d21c3543f59e981ad8caa505177f6f995d27dd2fc0ed2e",
      "occurred_at": "2026-08-25T00:01:06.988832836Z",
      "payload_digest": "sha256:1da929872d7c566a5426c2170011fcd0173f5ee23f24987bd86dae5458ef1e37"
    },
    {
      "seq": 39,
      "kind": "UsageCommitted",
      "subject": "work:c72c1d0d5e83480b41f7d31c4613bc630efe7c38eb4e41cb5af627a93fae25cc",
      "occurred_at": "2026-08-25T00:01:06.999462443Z",
      "payload_digest": "sha256:96d43da28df8e28f71cb92e8c012ab2d64dc4f7348e872a5d42038e6ace821f9"
    },
    {
      "seq": 40,
      "kind": "EpisodeCompleted",
      "subject": "episode:f67a57c5a2e6dfe1d0d21c3543f59e981ad8caa505177f6f995d27dd2fc0ed2e",
      "occurred_at": "2026-08-25T00:01:06.999997421Z",
      "payload_digest": "sha256:a55d9a22ec39e4e2a83710e271af50c92110f9637a7adc18f68c0cd466cd7e53"
    },
    {
      "seq": 41,
      "kind": "ObservationRecorded",
      "subject": "observation:29254b9749017475a2fcaef92cc73f3912e1045b3e761dbec2c100c182c08554",
      "occurred_at": "2026-08-25T00:01:07.000507805Z",
      "payload_digest": "sha256:0fa633d59698492794fc166f5241042b9a676ed9ee467b20874fc02ce4e6257c"
    },
    {
      "seq": 42,
      "kind": "ObservationRecorded",
      "subject": "observation:a56c56da3629a81ce8f2fe0b9a97e49286071cdfb3f73ea8bc2f4cf7b87c1147",
      "occurred_at": "2026-08-25T00:01:07.0010273Z",
      "payload_digest": "sha256:b85ee58b935944200681cd5bbdb919bea0b16e13578eb38fc8f137e24827c714"
    },
    {
      "seq": 43,
      "kind": "ObservationRecorded",
      "subject": "observation:1ac1149c1407be822d214cda9a69c8e4ce1e01a5a0171ce551c66736a3b929f0",
      "occurred_at": "2026-08-25T00:01:07.001550614Z",
      "payload_digest": "sha256:487fa8f3d316f61400cd7da9d1716c77cfde0580426fac4cf5e5dcbc2fee901e"
    },
    {
      "seq": 44,
      "kind": "EpisodeLeaseGranted",
      "subject": "episode:b4c54b1edb5ecec34d55170212c72e95f0a2f401224b1156f06581b69bd62ac1",
      "occurred_at": "2026-08-25T00:01:07.003014733Z",
      "payload_digest": "sha256:f3d04c3017221fbf52d2eaa0c6eb263fdfbaa64dd67d7cfa72d40e1130457d1d"
    },
    {
      "seq": 45,
      "kind": "EpisodeAttemptStarted",
      "subject": "episode:b4c54b1edb5ecec34d55170212c72e95f0a2f401224b1156f06581b69bd62ac1",
      "occurred_at": "2026-08-25T00:01:07.004227953Z",
      "payload_digest": "sha256:db518dc381c1b9e0a6557f68a730f882a4a59926d67e5c22b82cc0f432d5f4b6"
    },
    {
      "seq": 46,
      "kind": "UsageCommitted",
      "subject": "work:ebaec2304a03a8f44079ea1f6020a1ecc94e35fdfb5215dd6fd7dfadb311fed5",
      "occurred_at": "2026-08-25T00:01:07.015686335Z",
      "payload_digest": "sha256:40e09e5fe942b76244c5f1026e40cf92cf083d6cd913cfb7c56e60fa497789ca"
    },
    {
      "seq": 47,
      "kind": "EpisodeCompleted",
      "subject": "episode:b4c54b1edb5ecec34d55170212c72e95f0a2f401224b1156f06581b69bd62ac1",
      "occurred_at": "2026-08-25T00:01:07.016698187Z",
      "payload_digest": "sha256:24384f2f234059b61d7c1acf736603f952bb436a450116679fe25ed4002ed576"
    },
    {
      "seq": 48,
      "kind": "ObservationRecorded",
      "subject": "observation:d9e4ce4815a4af1fd5aba952a3c21cd61fe748063ec03cb0671b99b133447cc2",
      "occurred_at": "2026-08-25T00:01:07.017711311Z",
      "payload_digest": "sha256:db15a7f53b688ed31b0074cc4d89f398128b21cd64e280d8a7e3c30b573f813c"
    },
    {
      "seq": 49,
      "kind": "ObservationRecorded",
      "subject": "observation:0806a1ee49bf75fd1fd55c0ed43644aecd7c722b6ee153d80a8b687c859a07e3",
      "occurred_at": "2026-08-25T00:01:07.018760339Z",
      "payload_digest": "sha256:b5f0965d84101ca637e26c519ae0299a20f920157edec53883f7c86e42b1b516"
    },
    {
      "seq": 50,
      "kind": "ObservationRecorded",
      "subject": "observation:7c85b000db482cec995bc746e458d3041af244d936a872a35d8a45473d058076",
      "occurred_at": "2026-08-25T00:01:07.019782479Z",
      "payload_digest": "sha256:6d23af8360fecb6fd7d6d84bccec9d2c1c59e19b39d4ef8ca42e6503f2232da5"
    },
    {
      "seq": 51,
      "kind": "EpisodeLeaseGranted",
      "subject": "episode:9e32c98093ea391c1c2e7ba89dda29f29ebad2cc1963175106c592b84d9948be",
      "occurred_at": "2026-08-25T00:01:07.022365389Z",
      "payload_digest": "sha256:cdfafd94d78f38757da749ed7c1644fc8834a59483a97c24fcd9541f2d5005d2"
    },
    {
      "seq": 52,
      "kind": "EpisodeAttemptStarted",
      "subject": "episode:9e32c98093ea391c1c2e7ba89dda29f29ebad2cc1963175106c592b84d9948be",
      "occurred_at": "2026-08-25T00:01:07.024124432Z",
      "payload_digest": "sha256:151a8ecfd8f98d36776bb883d61e387677beb6beecb3b492ec8fb7aa85d77f36"
    },
    {
      "seq": 53,
      "kind": "UsageCommitted",
      "subject": "work:c2a260ab117fc8013bf6650143ff63be8063d262ea31040b3f181b0eafe1f897",
      "occurred_at": "2026-08-25T00:01:07.036500274Z",
      "payload_digest": "sha256:7200463f0f0cbbca5bcf921ddaba309f842f94ed9f43db85296ca5247829dd96"
    },
    {
      "seq": 54,
      "kind": "EpisodeCompleted",
      "subject": "episode:9e32c98093ea391c1c2e7ba89dda29f29ebad2cc1963175106c592b84d9948be",
      "occurred_at": "2026-08-25T00:01:07.037200699Z",
      "payload_digest": "sha256:cb56b0df786e85a4414558be94a9fe24a621fd8acac65b39f7165e78e8a41ff8"
    },
    {
      "seq": 55,
      "kind": "ObservationRecorded",
      "subject": "observation:ddb1f28872f4439d24a0abcf8c4ce681ffb27a5801ff2d78e0497926319bba7a",
      "occurred_at": "2026-08-25T00:01:07.037901162Z",
      "payload_digest": "sha256:3df63c6c1e09958e18aaf63aee841fa9eb10181312629cf995fa1e3ef613bf45"
    },
    {
      "seq": 56,
      "kind": "ObservationRecorded",
      "subject": "observation:14312eafdd17ab11f423f9573117536050bedf653fed45fb1944ae8a16b16467",
      "occurred_at": "2026-08-25T00:01:07.03857853Z",
      "payload_digest": "sha256:3ff2251c4a5d9244a6c33d11d011fd1fe17fe76d0a8ae6f7ca739cab9eb1834e"
    },
    {
      "seq": 57,
      "kind": "ObservationRecorded",
      "subject": "observation:689682f1b615f838086c15a8401abeb4d1ff4cae3f63ace2a418add283cdc5e0",
      "occurred_at": "2026-08-25T00:01:07.039264834Z",
      "payload_digest": "sha256:1c564671055581eb21263011b9903ff278a1a5ffe179c589e381f15a3683b49f"
    },
    {
      "seq": 58,
      "kind": "EpisodeLeaseGranted",
      "subject": "episode:ed10febc1658b5170f8977165fcaa18254154326c01fe48b68b314ee4ac62fb5",
      "occurred_at": "2026-08-25T00:01:07.041115493Z",
      "payload_digest": "sha256:dcb8a0a8286f3b443ced9301127a1f5042c54883308e2403ee9d9fdc54191f17"
    },
    {
      "seq": 59,
      "kind": "EpisodeAttemptStarted",
      "subject": "episode:ed10febc1658b5170f8977165fcaa18254154326c01fe48b68b314ee4ac62fb5",
      "occurred_at": "2026-08-25T00:01:07.042729039Z",
      "payload_digest": "sha256:35ca8f14b4b54017fe07fcc758b5afda9a45c8482c4723ba027eaeb7406737ce"
    },
    {
      "seq": 60,
      "kind": "UsageCommitted",
      "subject": "work:5c5b85f6a813204a583d51dca202027245ba65e20c80af913cb437a1ad113f3e",
      "occurred_at": "2026-08-25T00:01:07.053304976Z",
      "payload_digest": "sha256:afaed993cdf62309c48da3b17b1a561d3e6e149db6cb6ff0de22e795a588af9b"
    },
    {
      "seq": 61,
      "kind": "EpisodeCompleted",
      "subject": "episode:ed10febc1658b5170f8977165fcaa18254154326c01fe48b68b314ee4ac62fb5",
      "occurred_at": "2026-08-25T00:01:07.054253362Z",
      "payload_digest": "sha256:31329810d06c528ec0f49c80fa584217a02cc36f8caa84e6cb91ae117f69ef5c"
    },
    {
      "seq": 62,
      "kind": "ObservationRecorded",
      "subject": "observation:0183429000d458a4ea5ed6270203c1d7d2ccbd3c69076167e65b870469d0d9fc",
      "occurred_at": "2026-08-25T00:01:07.05503021Z",
      "payload_digest": "sha256:6076b6e2f01329479dfcbd882f1a9695839f3cec6b2360700f9ffbd27877c3fe"
    },
    {
      "seq": 63,
      "kind": "ObservationRecorded",
      "subject": "observation:21a92dfce7e3d4b66d3683c018c920d0b6bb2a5e57f77acb5c1aec8a3de0f383",
      "occurred_at": "2026-08-25T00:01:07.055772337Z",
      "payload_digest": "sha256:906997c704fcb7627dd0b10220abb885b265c046389826c194e45cd11df80406"
    },
    {
      "seq": 64,
      "kind": "ObservationRecorded",
      "subject": "observation:294775bc7c7f166589fddb8fb6ea8b831ed8a1f164be92cb5eb686c8a52dd2af",
      "occurred_at": "2026-08-25T00:01:07.056512659Z",
      "payload_digest": "sha256:618852fd9e8c41f87ba0e64d2d4d8c4f964e61e86c6a22f5b4f0dbf113568f47"
    },
    {
      "seq": 65,
      "kind": "EpisodeLeaseGranted",
      "subject": "episode:a6eb91f286136e1eebb6b9c193c3bfd31377bd2d77c9ca2b14e7b93a028285f5",
      "occurred_at": "2026-08-25T00:01:07.058428702Z",
      "payload_digest": "sha256:4a325f55836c85af5811ffa26c43ab48ba02e5762dec38cb321c2a702b500b83"
    },
    {
      "seq": 66,
      "kind": "EpisodeAttemptStarted",
      "subject": "episode:a6eb91f286136e1eebb6b9c193c3bfd31377bd2d77c9ca2b14e7b93a028285f5",
      "occurred_at": "2026-08-25T00:01:07.060061032Z",
      "payload_digest": "sha256:83b2794363a50260bf37eaa594f9cd44098d81789f2dd9163b888d1795f43afc"
    },
    {
      "seq": 67,
      "kind": "UsageCommitted",
      "subject": "work:adabf70aa0e7f84c4465f40eb9cab5ca05528879c1989147ffac336c374c78a2",
      "occurred_at": "2026-08-25T00:01:07.07147811Z",
      "payload_digest": "sha256:4738ef700646d7291cb12a3f4ed479c14ead24c97e1bec60cb7aba1ceeafe9da"
    },
    {
      "seq": 68,
      "kind": "EpisodeCompleted",
      "subject": "episode:a6eb91f286136e1eebb6b9c193c3bfd31377bd2d77c9ca2b14e7b93a028285f5",
      "occurred_at": "2026-08-25T00:01:07.072321571Z",
      "payload_digest": "sha256:7ddaf7f19442bb306ba81e813ae964b43d6a71d87259597e8162b215e1dffb00"
    },
    {
      "seq": 69,
      "kind": "ObservationRecorded",
      "subject": "observation:7d7dda772f598dcabf6383b4528f61872a47100013742ba96a6837ae2cfdeee6",
      "occurred_at": "2026-08-25T00:01:07.073165059Z",
      "payload_digest": "sha256:69f91003c8ee9404e52bc42ffdfcc94d466f4a59a136fe70b3547702136aed3e"
    },
    {
      "seq": 70,
      "kind": "ObservationRecorded",
      "subject": "observation:ef55a878cff198a34467dc0c827e1ced0a84920802757634417d6fac47daef50",
      "occurred_at": "2026-08-25T00:01:07.073994973Z",
      "payload_digest": "sha256:21fde09c771884b58c5b1766947cfaea0259ab6458de24e0235548212c2debf9"
    },
    {
      "seq": 71,
      "kind": "ObservationRecorded",
      "subject": "observation:15122abda8fd14e86556f6643cb9f55abbb523599b7a8dacdbeb5692c0a1ffe1",
      "occurred_at": "2026-08-25T00:01:07.07487526Z",
      "payload_digest": "sha256:ae4040ed680e90bfaf5c8cc09adfa2c06269d49da14e19339905b471363c8ece"
    },
    {
      "seq": 72,
      "kind": "EpisodeLeaseGranted",
      "subject": "episode:1b8d31e186897e58a3aa00fdf51a4893a2c8f867bd04e5af0d0d0bf1e81c5f7d",
      "occurred_at": "2026-08-25T00:01:07.07672133Z",
      "payload_digest": "sha256:cf495ac5ceb083fb9c1e97e2f1e094e3f5a2bbfd5ca3bedd1c88e0bd714e96b5"
    },
    {
      "seq": 73,
      "kind": "EpisodeAttemptStarted",
      "subject": "episode:1b8d31e186897e58a3aa00fdf51a4893a2c8f867bd04e5af0d0d0bf1e81c5f7d",
      "occurred_at": "2026-08-25T00:01:07.078289898Z",
      "payload_digest": "sha256:71d0959f222b79041d61cd20bfabb13f74375d2c58e39cddaa2080ea76991527"
    },
    {
      "seq": 74,
      "kind": "UsageCommitted",
      "subject": "work:2c0685b50c91adab89e47180f469b660b712c00cbef4d3d734d1e448a064862b",
      "occurred_at": "2026-08-25T00:01:07.088777364Z",
      "payload_digest": "sha256:1fb396908bd626aa7ec521e55d3a6163e3ae45fae9f7238ea9910718f8832144"
    },
    {
      "seq": 75,
      "kind": "EpisodeCompleted",
      "subject": "episode:1b8d31e186897e58a3aa00fdf51a4893a2c8f867bd04e5af0d0d0bf1e81c5f7d",
      "occurred_at": "2026-08-25T00:01:07.089662667Z",
      "payload_digest": "sha256:d4f105709a46c5af8b2b5346f9704850a6ac6627f36a21d454144cce21839443"
    },
    {
      "seq": 76,
      "kind": "ObservationRecorded",
      "subject": "observation:f028a2a3c60f94bca694e541b2bde50937025587221d8efd6454e6be0a0e2020",
      "occurred_at": "2026-08-25T00:01:07.090549967Z",
      "payload_digest": "sha256:525abb9ab42d2a3f5d792b1730ead894f3934a45f47610c7f662face69d973b2"
    },
    {
      "seq": 77,
      "kind": "ObservationRecorded",
      "subject": "observation:48601c849902f188b0b4f9cceff3c6ee2fa891bdd357a04f04b1fbffef6ee3ef",
      "occurred_at": "2026-08-25T00:01:07.09144446Z",
      "payload_digest": "sha256:1ceac170728ae0dd15064af53e6b77027b9349658f409ebc29308eb89e497ba2"
    },
    {
      "seq": 78,
      "kind": "ObservationRecorded",
      "subject": "observation:d7f757d1c2127a8887e9489483d08e57e05b55b6d1cbbc18d602beae8da683d5",
      "occurred_at": "2026-08-25T00:01:07.092360202Z",
      "payload_digest": "sha256:108507c134c0c94f897ae111d70256d2f85d86681d4521effcb58fdd45b7ddb6"
    },
    {
      "seq": 79,
      "kind": "EstimateRecorded",
      "subject": "trial:3347f0684e4c8e061316284058e903cfc0e6618231294a474c0f688c5d3e6619",
      "occurred_at": "2026-08-25T00:01:07.101994444Z",
      "payload_digest": "sha256:3a6f7835849d0ef402e4a41e1c6335bde4b7c2b28bb679f7f2b10037460debc2"
    },
    {
      "seq": 80,
      "kind": "DecisionRecorded",
      "subject": "candidate:7352fad2bfeac8caa485885fc767bb8cd031c6b960ce679a3658d228e64a44e9",
      "occurred_at": "2026-08-25T00:01:07.10426413Z",
      "payload_digest": "sha256:6c227644b60239bea4d212b1e8474f527c283165406be3e1dde084048b981632"
    },
    {
      "seq": 81,
      "kind": "CampaignCompleted",
      "occurred_at": "2026-08-25T00:01:07.106080194Z",
      "payload_digest": "sha256:c5f5f221b4fece1ce535cb8a212137a2c5e26c9e7b5c86d5dd38372648a50cad"
    }
  ]
}
```

## Verification

```json
{
  "campaign": "campaign:b200ae2487f1eed38eabab8f24cc3cfd",
  "events": 81,
  "journal": "verified",
  "unique_event_payloads": 80
}
```
