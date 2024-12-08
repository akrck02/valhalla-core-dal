package org.akrck02.valhalla.core.dal.database

import com.mongodb.client.model.Filters
import com.mongodb.kotlin.client.coroutine.MongoClient
import com.mongodb.kotlin.client.coroutine.MongoDatabase
import org.akrck02.valhalla.core.dal.configuration.DatabaseConnectionConfiguration
import org.akrck02.valhalla.core.sdk.model.exception.ServiceException
import org.bson.conversions.Bson
import org.bson.types.ObjectId

class Mongo : Database {

    var client: MongoClient? = null

    /**
     * {inheritDoc}
     */
    override fun connect(configuration: DatabaseConnectionConfiguration) {
        client = MongoClient.create(
            connectionString = getConnectionString(configuration)
        )
    }

    fun getDatabase(database: Databases): MongoDatabase {
        val database = client?.getDatabase(databaseName = Databases.ValhallaTest.id)
        database ?: throw ServiceException(message = "Cannot connect to database")
        return database
    }

    /**
     * Get the connection string based in th parameters
     * @param configuration The connection configuration
     */
    private fun getConnectionString(configuration: DatabaseConnectionConfiguration) =
        "mongodb://${configuration.user}:${configuration.password}@${configuration.host}/?retryWrites=true&w=majority"


}

fun mongoIdEquals(id: String?): Bson {
    return Filters.eq("_id", ObjectId(id))
}