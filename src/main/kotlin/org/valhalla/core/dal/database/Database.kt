package org.valhalla.core.dal.database

import org.valhalla.core.dal.configuration.DatabaseConnectionConfiguration

interface Database {

    /**
     * Connect to the database
     * @param configuration The connection configuration
     */
    fun connect(configuration: DatabaseConnectionConfiguration)

}