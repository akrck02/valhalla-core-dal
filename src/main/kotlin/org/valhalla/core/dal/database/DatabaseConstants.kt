package org.valhalla.core.dal.database

/**
 * This enumeration represents the different databases used by the service
 */
enum class DatabaseIdentifier(val id: String) {
    Valhalla("valhalla"),
    ValhallaTest("valhalla-test")
}

/**
 * This enumeration represents the different collections or tables used by the service
 */
enum class DatabaseCollections(val id: String) {
    Users("users")
}