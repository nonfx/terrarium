# This file is maintained automatically by "terraform init".
# Manual edits may be lost in future updates.

provider "registry.terraform.io/hashicorp/aws" {
<<<<<<< HEAD
  version     = "5.3.0"
  constraints = ">= 3.29.0, >= 3.72.0, >= 3.74.0, >= 4.0.0, >= 4.28.0, >= 4.35.0"
  hashes = [
    "h1:iLQKGkHnita0bqbWUPulvgQaMY92azVG+pYPGfd66Jg=",
    "zh:001814dcf6b2329de5e2c9223c4f1e95a0f60d6670046015419053b03b3c0712",
    "zh:3c511a91f53076c3a1117526bee0880b339261f1eb3feecd7854771bfef7890d",
    "zh:3e6c19e048f06051c9296c7a3236946f37431ce0d84f843585c5f3e8504759d3",
    "zh:476a3d918782a479166f33418192b522698e39702e8a0aec823682d3ee3082f1",
    "zh:5dd0d3bff7a7acabeed600dfbbef797e189c4877f65e4b4ed572cb33e454f602",
    "zh:6627f95a41e30c01b7f7c9e3db1cccba056c5257c36cccfaa0898d526211add2",
    "zh:663023a4244cf7f7df2b08ab204922f7902eefe9a7b51a2c2def1a7dafe6f55f",
    "zh:79cb8a22a131b7d2beb331d8443207eed10fdb4b09655048960bd5d59c8bbf3a",
    "zh:8c2275a0954042cfc44843a6045543744e08bd8cad487f0bc9162cf92a9bcdcc",
=======
  version     = "4.35.0"
  constraints = ">= 3.0.0, >= 3.29.0, >= 3.72.0, >= 3.74.0, >= 4.0.0, >= 4.35.0, 4.35.0, < 5.0.0"
  hashes = [
    "h1:YWGliEq8S7vVrR+I/lwr9GcyVctB1n9/Qz7eElKrXbg=",
    "zh:045c9c113311a358e6f311a6d7c67f4a18a53d6468a5dbe4ad4d1c5a3cf089cf",
    "zh:43ee43aca5a5377e3b55463c19ab497e24a3653c233151214d1907ff3d7ae749",
    "zh:5834362e4a402bb2682de4166c340fdc88c910d393c1753c613a526685279083",
    "zh:64a0066e1893077d70aaa13f2ab7a9e3a5bc676767daa4036088e28c799a5b88",
    "zh:690cbc4cfad5f74899bd0695896ecd1e9cb3dd362dfcae13701eb5e955409372",
    "zh:82ebbd737671bf8f4ed85183c4a37115ae7fc6aa9a6213e30509a4f806e593a0",
    "zh:8b9a92114b09eadd594f8f39edadaa103e640d57a10df3b7283a875d76faf2e4",
>>>>>>> 2552b7c (VAN-4199 Improve resource-to-resource link detection)
    "zh:9b12af85486a96aedd8d7984b0ff811a4b42e3d88dad1a3fb4c0b580d04fa425",
    "zh:ad08ae20b9402461af863772a9e4ff5677e14f3fc86d5b148bd4faaaa361f601",
    "zh:b8b7bd15fc1842aeedc2e5eab03b8357cdb2b9fe3e67dd82ae240be3081bf637",
    "zh:bdb3858c4c632aad8d5c4bff063f3afb18de51cec3167b3496d5bc5856915301",
    "zh:f354a433ec8095b06c2701725411ffb73a20ef9b1aa325434e1bb575b5c86d52",
    "zh:f47e1342883d599f4675dcfdeb9707cdfcfaf53c677f93fd5c410580d4dece13",
  ]
}

provider "registry.terraform.io/hashicorp/cloudinit" {
  version     = "2.3.2"
  constraints = ">= 2.0.0"
  hashes = [
    "h1:Ar/DAbZQ9Nsj0BrqX6camrEE6U+Yq4E87DCNVqxqx8k=",
    "zh:2487e498736ed90f53de8f66fe2b8c05665b9f8ff1506f751c5ee227c7f457d1",
    "zh:3d8627d142942336cf65eea6eb6403692f47e9072ff3fa11c3f774a3b93130b3",
    "zh:434b643054aeafb5df28d5529b72acc20c6f5ded24decad73b98657af2b53f4f",
    "zh:436aa6c2b07d82aa6a9dd746a3e3a627f72787c27c80552ceda6dc52d01f4b6f",
    "zh:458274c5aabe65ef4dbd61d43ce759287788e35a2da004e796373f88edcaa422",
    "zh:54bc70fa6fb7da33292ae4d9ceef5398d637c7373e729ed4fce59bd7b8d67372",
    "zh:78d5eefdd9e494defcb3c68d282b8f96630502cac21d1ea161f53cfe9bb483b3",
    "zh:893ba267e18749c1a956b69be569f0d7bc043a49c3a0eb4d0d09a8e8b2ca3136",
    "zh:95493b7517bce116f75cdd4c63b7c82a9d0d48ec2ef2f5eb836d262ef96d0aa7",
    "zh:9ae21ab393be52e3e84e5cce0ef20e690d21f6c10ade7d9d9d22b39851bfeddc",
    "zh:cc3b01ac2472e6d59358d54d5e4945032efbc8008739a6d4946ca1b621a16040",
    "zh:f23bfe9758f06a1ec10ea3a81c9deedf3a7b42963568997d84a5153f35c5839a",
  ]
}

provider "registry.terraform.io/hashicorp/google" {
  version     = "3.90.1"
  constraints = ">= 3.22.0, >= 3.28.0, < 4.0.0"
  hashes = [
    "h1:9TYwyR4R4dIop7wV2lvvYZHw9RUVd/YRWR+9jjXpyfw=",
    "zh:07aabc8e46a5a2b29932e10677b23d4ce9d9a25f22ab61d3307a6b0e7998c84e",
    "zh:0b63cd9534a98ed0fee794da495833046ad5319bd2da3102e21a941b7e2b857e",
    "zh:17f815d57e1426edf8818323ab8e1022c8ec60dce0ced89a3b8e5dde5a95b3cc",
    "zh:37855eae3542f2ebc6416984b124533d00299e0e01dcd7d2bc2205469cb9eceb",
    "zh:579aa32a8e3fa317ddbd28c99a6449ae8864a5b7d10247bca6496f399cb36701",
    "zh:703f71e0231cfe7a025c61db361d928189adba1d4fad2fe77f783dc73c8afe30",
    "zh:afcd80c31cb1ed75ce6813269618e01ab29af68dae7aae1c51521c13acdaa678",
    "zh:b21302f65a0d37045216912695d1ef718a1fe1732c30dc5654891fe2519b8e4e",
    "zh:b69d0c8a74c2cd6233681db37e01aaaf1a6fb6bb24c83f7715bd2b456083e29d",
    "zh:d4fb305816b143cb26c1827c79e56651347fd41809a57184e4807fb3f804f510",
    "zh:fa24173ef9524bdfa1c5cada5188489554b08374f9519fe545f3fc1d3a9d9d4f",
  ]
}

provider "registry.terraform.io/hashicorp/google-beta" {
  version = "4.69.1"
  hashes = [
<<<<<<< HEAD
    "h1:P0FZVE06mpY+7LSraLChJ1VMv0njSFw3I1Kymhr2bJ4=",
=======
    "h1:itYpL6typEinqgQLaBN67SCuDuK0Zm0/o/w6vDEuczY=",
>>>>>>> 2552b7c (VAN-4199 Improve resource-to-resource link detection)
    "zh:0e5deb489aa9c80fc039388c9d8a08e3e20b625e3cd5c6aac57c67d9e56b5cbd",
    "zh:14ece88b3866f49ae1bc27f9e81ed1202aba3f69ec62123242643b099746fda7",
    "zh:1b7cefb2569f4063bd412073fa7444ecc0b96d6e8e065303fa57ed7c6287e328",
    "zh:282d9375665a04bd4d35eec605736680f86ed29ab2c548acbffea1a0557e1009",
    "zh:29125aae9e38355718ecf29d547d6557d4c3aa41f08e46c067995eac3f264fc3",
    "zh:30826b38fadd23e518451c63846547fe31bb3d846b87a3f3959ae73bc3352809",
    "zh:9ccdad51407ba2a67b9a93471550e6076744dfd02db15b588c6b0b6c1b1de885",
    "zh:a980488ec165fbdffc0e7e2b9822340564c800094d345bea527a17ac7ff7c677",
    "zh:b3aa35373e7a72f4ea4ef3983cfde4dfd479bff3cfb066f565fcbc942ecefbfb",
    "zh:c45f35df9dda13bc854813d030c215c3f75cffea026830a42c6166407d78d597",
    "zh:e54293c1cbf977e4fc99e7593ba6d046548c948ae8fe5b39d9b2221645c222a9",
    "zh:f569b65999264a9416862bca5cd2a6177d94ccb0424f3a4ef424428912b9cb3c",
  ]
}

provider "registry.terraform.io/hashicorp/helm" {
  version     = "2.1.2"
  constraints = "2.1.2"
  hashes = [
    "h1:UVuNjmuEM4ZVtItbh1QRGulkBWxDY929roxFQhEf9Ks=",
    "zh:09bd2b6f33a040c3fd59d82c9768b886b8c82163e31ec92dc1b747229d0548df",
    "zh:09f209fa57ad5d01f04c458f1719b42958ca5e0fc2eca63d9ec29f92c77a29f8",
    "zh:0bfc627539500ffb2a41a2f8a5ea7f6fb1d76367b11bbf9489b483b9e8dfff8f",
    "zh:0c0fef5587a5e927d15f9f4cc13cd0620b138238f9a422490fe9ea2bf086b61a",
    "zh:187f99648fad2b84d49cdd372f8f6cedbf06e13411b3f1ff66708f66852d7855",
    "zh:3d9ae08f8a99b19e80bd27708aecf592c28c92da66fd60189dfd7dce4d7da93c",
    "zh:60b767109362c616b2e6386bfb08581b03bc3e528920444e52b16743f5a180d6",
    "zh:729db42ed49d91c9b51eb602b9253e6ed6b3ab613c42deefc14996c9a8ee8ae4",
    "zh:8401f3bf6d69ce43eb14911823c7e5cbb273cf564508043cd04fb064c30a3e1a",
    "zh:91139b492ce1f41847017349ea49f9441b7cf70762c8d1c32a6a909e25ed10c1",
    "zh:98fca606a539510edc94dcad8069a321e6a42df90e483f58df03b305726d9220",
  ]
}

provider "registry.terraform.io/hashicorp/http" {
  version = "3.3.0"
  hashes = [
    "h1:QL/rtSlbi+F+ukbr/k4MahiO5lX4AiEu37p4kOV9ELk=",
    "zh:27d101f4c089d1e367bbbbb3f260fc7d52f63559a4424c08633e566863c951b2",
    "zh:37860671324229f52a7d82eea88a31fe24321297fd699d879de5b6cf6aae086c",
    "zh:4680716579e361298e4331ce0c92e38011fc41ed56bd55302c23b696b3b8c469",
    "zh:547cd2a407ca0d22307634d83ffc64cd4225f221baa09682b7a8c5a2429c34d8",
    "zh:61965698af75aad7482f2f593b75f15e4a4f6f0117b643c69f3da61f40b1a9c7",
    "zh:78d5eefdd9e494defcb3c68d282b8f96630502cac21d1ea161f53cfe9bb483b3",
    "zh:93f9e0f2244816cbb72197c733ada4214df691e4e6a84b8e340e43e43ab8a383",
    "zh:969aad70624d033c257c365cf75001d29fa7341b48d673cd7317205395b4791b",
    "zh:e9504018b1af992c041bda1e4a6f01db1f1cdb1a7df8055d1082049befbc4217",
    "zh:fa7f6af94e75c6fe21782c622ed387ae08ee3ffeaa0176f08d0b06bb61bb50f4",
    "zh:feda1d7cdae86bce829f82223f625b55c858a36d3aca1a762d7258798a25b476",
    "zh:ff1f3d8c53930aad2fde32d6328df7e7e5b5de36dd7c0682d15518993ab199ef",
  ]
}

provider "registry.terraform.io/hashicorp/kubernetes" {
  version     = "2.21.1"
  constraints = ">= 2.10.0"
  hashes = [
    "h1:2spGoBcGDQ/Csc23bddCfM21zyKx3PONoiqRgmuChnM=",
    "zh:156a437d7edd6813e9cb7bdff16ebce28cec08b07ba1b0f5e9cec029a217bc27",
    "zh:1a21c255d8099e303560e252579c54e99b5f24f2efde772c7e39502c62472605",
    "zh:27b2021f86e5eaf6b9ee7c77d7a9e32bc496e59dd0808fb15a5687879736acf6",
    "zh:31fa284c1c873a85c3b5cfc26cf7e7214d27b3b8ba7ea5134ab7d53800894c42",
    "zh:4be9cc1654e994229c0d598f4e07487fc8b513337de9719d79b45ce07fc4e123",
    "zh:5f684ed161f54213a1414ac71b3971a527c3a6bfbaaf687a7c8cc39dcd68c512",
    "zh:6d58f1832665c256afb68110c99c8112926406ae0b64dd5f250c2954fc26928e",
    "zh:9dadfa4a019d1e90decb1fab14278ee2dbefd42e8f58fe7fa567a9bf51b01e0e",
    "zh:a68ce7208a1ef4502528efb8ce9f774db56c421dcaccd3eb10ae68f1324a6963",
    "zh:acdd5b45a7e80bc9d254ad0c2f9cb4715104117425f0d22409685909a790a6dd",
    "zh:f569b65999264a9416862bca5cd2a6177d94ccb0424f3a4ef424428912b9cb3c",
    "zh:fb451e882118fe92e1cb2e60ac2d77592f5f7282b3608b878b5bdc38bbe4fd5b",
  ]
}

provider "registry.terraform.io/hashicorp/null" {
  version = "3.2.1"
  hashes = [
    "h1:tSj1mL6OQ8ILGqR2mDu7OYYYWf+hoir0pf9KAQ8IzO8=",
    "zh:58ed64389620cc7b82f01332e27723856422820cfd302e304b5f6c3436fb9840",
    "zh:62a5cc82c3b2ddef7ef3a6f2fedb7b9b3deff4ab7b414938b08e51d6e8be87cb",
    "zh:63cff4de03af983175a7e37e52d4bd89d990be256b16b5c7f919aff5ad485aa5",
    "zh:74cb22c6700e48486b7cabefa10b33b801dfcab56f1a6ac9b6624531f3d36ea3",
    "zh:78d5eefdd9e494defcb3c68d282b8f96630502cac21d1ea161f53cfe9bb483b3",
    "zh:79e553aff77f1cfa9012a2218b8238dd672ea5e1b2924775ac9ac24d2a75c238",
    "zh:a1e06ddda0b5ac48f7e7c7d59e1ab5a4073bbcf876c73c0299e4610ed53859dc",
    "zh:c37a97090f1a82222925d45d84483b2aa702ef7ab66532af6cbcfb567818b970",
    "zh:e4453fbebf90c53ca3323a92e7ca0f9961427d2f0ce0d2b65523cc04d5d999c2",
    "zh:e80a746921946d8b6761e77305b752ad188da60688cfd2059322875d363be5f5",
    "zh:fbdb892d9822ed0e4cb60f2fedbdbb556e4da0d88d3b942ae963ed6ff091e48f",
    "zh:fca01a623d90d0cad0843102f9b8b9fe0d3ff8244593bd817f126582b52dd694",
  ]
}

provider "registry.terraform.io/hashicorp/random" {
  version     = "3.5.1"
  constraints = ">= 3.1.0"
  hashes = [
    "h1:sZ7MTSD4FLekNN2wSNFGpM+5slfvpm5A/NLVZiB7CO0=",
    "zh:04e3fbd610cb52c1017d282531364b9c53ef72b6bc533acb2a90671957324a64",
    "zh:119197103301ebaf7efb91df8f0b6e0dd31e6ff943d231af35ee1831c599188d",
    "zh:4d2b219d09abf3b1bb4df93d399ed156cadd61f44ad3baf5cf2954df2fba0831",
    "zh:6130bdde527587bbe2dcaa7150363e96dbc5250ea20154176d82bc69df5d4ce3",
    "zh:6cc326cd4000f724d3086ee05587e7710f032f94fc9af35e96a386a1c6f2214f",
    "zh:78d5eefdd9e494defcb3c68d282b8f96630502cac21d1ea161f53cfe9bb483b3",
    "zh:b6d88e1d28cf2dfa24e9fdcc3efc77adcdc1c3c3b5c7ce503a423efbdd6de57b",
    "zh:ba74c592622ecbcef9dc2a4d81ed321c4e44cddf7da799faa324da9bf52a22b2",
    "zh:c7c5cde98fe4ef1143bd1b3ec5dc04baf0d4cc3ca2c5c7d40d17c0e9b2076865",
    "zh:dac4bad52c940cd0dfc27893507c1e92393846b024c5a9db159a93c534a3da03",
    "zh:de8febe2a2acd9ac454b844a4106ed295ae9520ef54dc8ed2faf29f12716b602",
    "zh:eab0d0495e7e711cca367f7d4df6e322e6c562fc52151ec931176115b83ed014",
  ]
}

provider "registry.terraform.io/hashicorp/template" {
  version = "2.2.0"
  hashes = [
    "h1:0wlehNaxBX7GJQnPfQwTNvvAf38Jm0Nv7ssKGMaG6Og=",
    "zh:01702196f0a0492ec07917db7aaa595843d8f171dc195f4c988d2ffca2a06386",
    "zh:09aae3da826ba3d7df69efeb25d146a1de0d03e951d35019a0f80e4f58c89b53",
    "zh:09ba83c0625b6fe0a954da6fbd0c355ac0b7f07f86c91a2a97849140fea49603",
    "zh:0e3a6c8e16f17f19010accd0844187d524580d9fdb0731f675ffcf4afba03d16",
    "zh:45f2c594b6f2f34ea663704cc72048b212fe7d16fb4cfd959365fa997228a776",
    "zh:77ea3e5a0446784d77114b5e851c970a3dde1e08fa6de38210b8385d7605d451",
    "zh:8a154388f3708e3df5a69122a23bdfaf760a523788a5081976b3d5616f7d30ae",
    "zh:992843002f2db5a11e626b3fc23dc0c87ad3729b3b3cff08e32ffb3df97edbde",
    "zh:ad906f4cebd3ec5e43d5cd6dc8f4c5c9cc3b33d2243c89c5fc18f97f7277b51d",
    "zh:c979425ddb256511137ecd093e23283234da0154b7fa8b21c2687182d9aea8b2",
  ]
}

provider "registry.terraform.io/hashicorp/tls" {
  version     = "4.0.4"
  constraints = ">= 3.0.0"
  hashes = [
    "h1:Wd3RqmQW60k2QWPN4sK5CtjGuO1d+CRNXgC+D4rKtXc=",
    "zh:23671ed83e1fcf79745534841e10291bbf34046b27d6e68a5d0aab77206f4a55",
    "zh:45292421211ffd9e8e3eb3655677700e3c5047f71d8f7650d2ce30242335f848",
    "zh:59fedb519f4433c0fdb1d58b27c210b27415fddd0cd73c5312530b4309c088be",
    "zh:5a8eec2409a9ff7cd0758a9d818c74bcba92a240e6c5e54b99df68fff312bbd5",
    "zh:5e6a4b39f3171f53292ab88058a59e64825f2b842760a4869e64dc1dc093d1fe",
    "zh:810547d0bf9311d21c81cc306126d3547e7bd3f194fc295836acf164b9f8424e",
    "zh:824a5f3617624243bed0259d7dd37d76017097dc3193dac669be342b90b2ab48",
    "zh:9361ccc7048be5dcbc2fafe2d8216939765b3160bd52734f7a9fd917a39ecbd8",
    "zh:aa02ea625aaf672e649296bce7580f62d724268189fe9ad7c1b36bb0fa12fa60",
    "zh:c71b4cd40d6ec7815dfeefd57d88bc592c0c42f5e5858dcc88245d371b4b8b1e",
    "zh:dabcd52f36b43d250a3d71ad7abfa07b5622c69068d989e60b79b2bb4f220316",
    "zh:f569b65999264a9416862bca5cd2a6177d94ccb0424f3a4ef424428912b9cb3c",
  ]
}
