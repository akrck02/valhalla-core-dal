package org.valhalla.core.dal.service.user.auth

import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.Test
import org.valhalla.core.dal.service.user.UserDataAccess
import org.valhalla.core.dal.tool.BaseDataAccessTest
import org.valhalla.core.sdk.repository.UserRepository

class UserAuthTest : BaseDataAccessTest() {

    companion object {
        lateinit var userRepository: UserRepository
    }

    @BeforeEach
    fun resetDatabaseImpl() {
        userRepository = UserDataAccess(database!!)
    }


    @Test
    fun login() {
        // TODO implement this
    }

    @Test
    fun authLogin() {
        // TODO implement this
    }

    @Test
    fun validateAccount() {
        // TODO implement this
    }
}