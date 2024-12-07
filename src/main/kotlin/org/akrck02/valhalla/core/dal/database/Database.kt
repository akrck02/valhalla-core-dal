package org.akrck02.valhalla.core.dal.database

import org.akrck02.valhalla.core.dal.configuration.DatabaseConnectionConfiguration

interface Database {

    /**
     * Connect to the database
     * @param configuration The connection configuration
     */
    fun connect(configuration: DatabaseConnectionConfiguration)
    
}