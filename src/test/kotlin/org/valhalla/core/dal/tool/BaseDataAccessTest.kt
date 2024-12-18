package org.valhalla.core.dal.tool

//import com.mongodb.kotlin.client.coroutine.MongoClient
//import com.mongodb.kotlin.client.coroutine.MongoDatabase
//import dal.configuration.getCurrentDatabaseConfiguration
//import dal.database.DatabaseIdentifier
//import dal.database.Mongo
import com.mongodb.kotlin.client.coroutine.MongoClient
import com.mongodb.kotlin.client.coroutine.MongoDatabase
import kotlinx.coroutines.runBlocking
import org.junit.jupiter.api.BeforeAll
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.TestInfo
import org.valhalla.core.dal.configuration.getCurrentDatabaseConfiguration
import org.valhalla.core.dal.database.DatabaseIdentifier
import org.valhalla.core.dal.database.Mongo
import org.valhalla.core.sdk.model.exception.ServiceException

/** Base data access test that prepares test database connection. */
open class BaseDataAccessTest {
    companion object {

        var client: MongoClient? = null
        var database: MongoDatabase? = null

        @JvmStatic
        @BeforeAll
        fun connectDatabase() {
            client = Mongo().also { it.connect(getCurrentDatabaseConfiguration()) }.client
        }
    }


    /** Reset the test database before each test. */
    @BeforeEach
    fun resetDatabase(info: TestInfo) = runBlocking {
        println("------------------------------------------")
        println("  ${info.displayName.removeSuffix("()").uppercase()}")
        println("------------------------------------------")
        println("Database: connecting to ${DatabaseIdentifier.ValhallaTest}")
        database = client?.getDatabase(databaseName = DatabaseIdentifier.ValhallaTest.id)
        database ?: throw ServiceException(message = "Cannot connect to database")
        println("Database: Cleaning ${DatabaseIdentifier.ValhallaTest}")
        database!!.drop()
    }

}