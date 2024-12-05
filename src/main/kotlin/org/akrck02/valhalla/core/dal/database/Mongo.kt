package org.akrck02.valhalla.core.dal.database

import com.mongodb.kotlin.client.coroutine.MongoClient
import org.akrck02.valhalla.core.dal.configuration.DatabaseConnectionConfiguration

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

    /**
     * Get the connection string based in th parameters
     * @param configuration The connection configuration
     */
    private fun getConnectionString(configuration: DatabaseConnectionConfiguration) =
        "mongodb://${configuration.user}:${configuration.password}@${configuration.host}/?retryWrites=true&w=majority"


}
