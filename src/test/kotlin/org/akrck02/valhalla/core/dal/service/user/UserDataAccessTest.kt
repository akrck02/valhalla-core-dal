package org.akrck02.valhalla.core.dal.service.user

import com.mongodb.kotlin.client.coroutine.MongoDatabase
import kotlinx.coroutines.launch
import kotlinx.coroutines.runBlocking
import org.akrck02.valhalla.core.dal.configuration.getCurrentDatabaseConfiguration
import org.akrck02.valhalla.core.dal.database.Databases
import org.akrck02.valhalla.core.dal.database.Mongo
import org.akrck02.valhalla.core.sdk.model.exception.ServiceException
import org.akrck02.valhalla.core.sdk.model.user.User
import org.junit.jupiter.api.BeforeEach
import kotlin.test.Test

class UserDataAccessTest {

    @BeforeEach
    fun resetDatabase() {
        runBlocking {
            val mongo = Mongo().also { it.connect(getCurrentDatabaseConfiguration()) }
            val database: MongoDatabase? = mongo.client?.getDatabase(databaseName = Databases.ValhallaTest.id)
            database ?: throw ServiceException(message = "Cannot connect to database")
        }
    }

    @Test
    fun userRegisterTest() {
        runBlocking {
            val mongo = Mongo().also { it.connect(getCurrentDatabaseConfiguration()) }
            val database: MongoDatabase? = mongo.client?.getDatabase(databaseName = Databases.ValhallaTest.id)
            launch {
                database ?: throw ServiceException(message = "Database not connected!")
                registerUser(
                    database,
                    User(
                        id = null,
                        username = "akrck01",
                        email = "akrck02@gmail.com",
                        password = "#PasswordisHereLoL#?"
                    )
                )
            }

        }
    }
}
