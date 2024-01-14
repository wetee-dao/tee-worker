package weteeguild

import types "wetee.app/worker/mint/chain/gen/types"

// See [`Pallet::guild_join`].
func MakeGuildJoinCall(daoId0 uint64, guildId1 uint64, who2 [32]byte) types.RuntimeCall {
	return types.RuntimeCall{
		IsWeteeGuild: true,
		AsWeteeGuildField0: &types.WeteeGuildPalletCall{
			IsGuildJoin:         true,
			AsGuildJoinDaoId0:   daoId0,
			AsGuildJoinGuildId1: guildId1,
			AsGuildJoinWho2:     who2,
		},
	}
}

// See [`Pallet::create_guild`].
func MakeCreateGuildCall(daoId0 uint64, name1 []byte, desc2 []byte, metaData3 []byte, creator4 [32]byte) types.RuntimeCall {
	return types.RuntimeCall{
		IsWeteeGuild: true,
		AsWeteeGuildField0: &types.WeteeGuildPalletCall{
			IsCreateGuild:          true,
			AsCreateGuildDaoId0:    daoId0,
			AsCreateGuildName1:     name1,
			AsCreateGuildDesc2:     desc2,
			AsCreateGuildMetaData3: metaData3,
			AsCreateGuildCreator4:  creator4,
		},
	}
}
