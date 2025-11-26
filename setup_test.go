package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type TestContext struct {
	context context.Context
	rdb     redis.UniversalClient
}

// TestContextInstance serves as a singleton instance of TestContext,
// ensuring that a consistent set of resources is available for all test cases.
var TestContextInstance *TestContext

func (tst *TestContext) Init() {
	// Create a new context for managing timeouts and cancellations in the Redis operations.
	// context.Background() provides an empty context that can be used as the root context.
	tst.context = context.Background()

	// Initialize the Redis client or cluster connection.
	// This step ensures that the Redis client (or cluster) is set up and ready for interactions within the test.
	// The _initRedis method configures the Redis instance to be used during the test for caching or data storage.
	tst._initRedis()

	// Store the initialized TestContext instance for future use.
	TestContextInstance = tst
}

// _initRedis initializes the Redis client and ensures the connection is ready for use in tests.
// It sets up the Redis Universal Client with the configuration provided in the test context.
func (tst *TestContext) _initRedis() {
	// Create a new Redis Universal Client using the provided configuration options.
	// The Addrs field specifies the Redis server addresses, and PoolSize determines the maximum number of connections.
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{Addrs: []string{"127.0.0.1:6379"}, PoolSize: 10, Password: "8!5s4n6$26WE!W"})

	// Test the connection to the Redis server by sending a ping command.
	// This ensures the Redis server is reachable and responsive before proceeding with tests.
	if err := rdb.Ping(tst.context).Err(); err != nil {
		// If the ping command fails, panic to immediately halt the test setup.
		// This indicates that the Redis server is not operational or unreachable.
		panic("failed redis ping")
	}

	// Assign the initialized Redis Universal Client to the redisCluster field in the test context.
	// This makes the Redis client available for use in test scenarios requiring Redis operations.
	tst.rdb = rdb
}

// GetRedis is a method on the TestContext struct that provides access to the Redis client instance.
// It returns a redis.UniversalClient, which is a unified interface for connecting to Redis.
// This method allows other parts of the code to retrieve the Redis client for performing Redis operations,
// enabling interaction with the Redis cluster managed by the TestContext.
func (tst *TestContext) GetRedis() redis.UniversalClient { return tst.rdb }

func (tst *TestContext) Stop() {
	// Close the Redis connection. This ensures that the Redis
	// client is properly terminated, releasing any open connections
	// to the Redis and preventing resource leaks.
	_ = tst.rdb.Close()
}

// TestMain is the entry point for executing tests in the package.
// It initializes the test context, sets up necessary resources, and ensures proper cleanup
// after the tests are executed. This function also provides logging for the start and end of the tests,
// ensuring that the execution flow is traceable. The use of os.Exit ensures the correct exit code
// is propagated based on the test results.
//func TestMain(m *testing.M) {
//	// Initialize a new instance of TestContext to manage resources and configuration for testing.
//	// This instance will provide shared functionality such as connections to external systems like Redis, ensuring a consistent test environment.
//	testContext := &TestContext{}
//	// Call the Init method on the TestContext instance to set up the test environment.
//	// The Init method initializes key components such as configuration, and
//	// external system connections, preparing the TestContext for use in tests.
//	testContext.Init()
//
//	// Ensure that resources managed by the TestContext are properly released after tests.
//	// The Stop method will clean up connections and other resources initialized during testing.
//	defer testContext.Stop()
//
//	log.Print("Init testing")
//
//	// Run the test suite using the provided test runner.
//	// The m.Run method executes all test functions in the current package and returns an exit code.
//	exitVal := m.Run()
//
//	log.Print("End testing")
//	// Exit the program with the appropriate exit code returned by the test runner.
//	// This ensures that the exit status reflects the success or failure of the tests.
//	os.Exit(exitVal)
//}
