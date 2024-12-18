package org.valhalla.core.dal.database

import com.mongodb.client.model.Filters
import com.mongodb.kotlin.client.coroutine.MongoClient
import com.mongodb.kotlin.client.coroutine.MongoDatabase
import org.bson.conversions.Bson
import org.bson.types.ObjectId
import org.valhalla.core.dal.configuration.DatabaseConnectionConfiguration
import org.valhalla.core.sdk.model.exception.ServiceException


@Suppress("unused")
class Mongo : Database {

    var client: MongoClient? = null
    override fun connect(configuration: DatabaseConnectionConfiguration) {
        client = MongoClient.create(
            connectionString = getConnectionString(configuration)
        )
    }

    /** Get database by [DatabaseIdentifier]. */
    fun getDatabase(database: DatabaseIdentifier): MongoDatabase {
        val db = client?.getDatabase(databaseName = database.id)
        return db ?: throw ServiceException(message = "Cannot connect to database")
    }

    /** Get the connection string based in the [DatabaseConnectionConfiguration] */
    private fun getConnectionString(configuration: DatabaseConnectionConfiguration) =
        "mongodb://${configuration.user}:${configuration.password}@${configuration.host}/?retryWrites=true&w=majority"


}

/** MongoDB [Filters] for the same id.  */
fun idEqualsFilter(id: String?): Bson {
    return Filters.eq("_id", ObjectId(id))
}