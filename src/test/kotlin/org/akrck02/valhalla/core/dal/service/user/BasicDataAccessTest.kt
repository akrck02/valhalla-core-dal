package org.akrck02.valhalla.core.dal.service.user

import com.mongodb.kotlin.client.coroutine.MongoClient
import com.mongodb.kotlin.client.coroutine.MongoDatabase
import kotlinx.coroutines.runBlocking
import org.akrck02.valhalla.core.dal.configuration.getCurrentDatabaseConfiguration
import org.akrck02.valhalla.core.dal.database.Databases
import org.akrck02.valhalla.core.dal.database.Mongo
import org.akrck02.valhalla.core.sdk.model.exception.ServiceException
import org.junit.jupiter.api.BeforeAll
import org.junit.jupiter.api.BeforeEach

open class BasicDataAccessTest {
    companion object {

        var client: MongoClient? = null
        var database: MongoDatabase? = null

        @JvmStatic
        @BeforeAll
        fun connectDatabase() {
            client = Mongo().also { it.connect(getCurrentDatabaseConfiguration()) }.client
        }
    }


    @BeforeEach
    fun resetDatabase() {
        runBlocking {
            database = client?.getDatabase(databaseName = Databases.ValhallaTest.id)
            database ?: throw ServiceException(message = "Cannot connect to database")
            database!!.drop()
        }
    }

}