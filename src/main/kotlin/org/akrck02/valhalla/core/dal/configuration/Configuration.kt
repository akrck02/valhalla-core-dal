package org.akrck02.valhalla.core.dal.configuration

import io.github.cdimascio.dotenv.dotenv

/**
 * This class represents the basic database connection configuration
 */
data class DatabaseConnectionConfiguration(
    val secret: String?,
    val host: String?,
    val user: String?,
    val password: String?
)

/**
 * Get the current database configuration
 */
fun getCurrentDatabaseConfiguration(): DatabaseConnectionConfiguration {
    val dotenv = dotenv { }

    return DatabaseConnectionConfiguration(
        user = dotenv["MONGO_USER"],
        password = dotenv["MONGO_PASSWORD"],
        host = dotenv["MONGO_SERVER"],
        secret = dotenv["SECRET"]
    )
}

