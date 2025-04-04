package balancer_test

import (
	"fmt"
	"testing"
	"time"

	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmtypes "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/osmosis-labs/osmosis/osmomath"
	"github.com/osmosis-labs/osmosis/v15/x/gamm/pool-models/balancer"
)

func createTestPool(t *testing.T, swapFee, exitFee sdk.Dec, poolAssets ...balancer.PoolAsset) *balancer.Pool {
	pool, err := balancer.NewBalancerPool(
		1,
		balancer.NewPoolParams(swapFee, exitFee, nil),
		poolAssets,
		"",
		time.Now(),
	)
	require.NoError(t, err)

	return &pool
}

func createTestContext(t *testing.T) sdk.Context {
	db := dbm.NewMemDB()
	logger := log.NewNopLogger()

	ms := rootmulti.NewStore(db, logger)

	return sdk.NewContext(ms, tmtypes.Header{}, false, logger)
}

func assertExpectedSharesErrRatio(t *testing.T, expectedShares, actualShares sdk.Int) {
	allowedErrRatioDec, err := sdk.NewDecFromStr(allowedErrRatio)
	require.NoError(t, err)

	errTolerance := osmomath.ErrTolerance{
		MultiplicativeTolerance: allowedErrRatioDec,
	}

	require.Equal(
		t,
		0,
		errTolerance.Compare(expectedShares, actualShares),
		fmt.Sprintf("expectedShares: %s, actualShares: %s", expectedShares.String(), actualShares.String()))
}

func assertExpectedLiquidity(t *testing.T, tokensJoined, liquidity sdk.Coins) {
	require.Equal(t, tokensJoined, liquidity)
}

// assertPoolStateNotModified asserts that sut (system under test) does not modify
// pool state.
func assertPoolStateNotModified(t *testing.T, pool *balancer.Pool, sut func()) {
	// We need to make sure that this method does not mutate state.
	oldPoolAssets := pool.GetAllPoolAssets()
	oldLiquidity := pool.GetTotalPoolLiquidity(sdk.Context{})
	oldShares := pool.GetTotalShares()

	sut()

	newPoolAssets := pool.GetAllPoolAssets()
	newLiquidity := pool.GetTotalPoolLiquidity(sdk.Context{})
	newShares := pool.GetTotalShares()

	require.Equal(t, oldPoolAssets, newPoolAssets)
	require.Equal(t, oldLiquidity, newLiquidity)
	require.Equal(t, oldShares, newShares)
}
