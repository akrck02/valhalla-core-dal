package org.akrck02.valhalla.core.dal.service.user.auth

import jdk.jshell.spi.ExecutionControl.NotImplementedException
import org.akrck02.valhalla.core.dal.service.user.UserDataAccess
import org.akrck02.valhalla.core.dal.tool.BaseDataAccessTest
import org.akrck02.valhalla.core.sdk.repository.UserRepository
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.Test

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
        throw NotImplementedException("not implemented yet!")
    }

    @Test
    fun authLogin() {
        throw NotImplementedException("not implemented yet!")
    }

    @Test
    fun validateAccount() {
        throw NotImplementedException("not implemented yet!")
    }

}