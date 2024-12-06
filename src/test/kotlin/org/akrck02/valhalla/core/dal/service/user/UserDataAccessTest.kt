package org.akrck02.valhalla.core.dal.service.user

import jdk.jshell.spi.ExecutionControl.NotImplementedException
import kotlinx.coroutines.runBlocking
import org.akrck02.valhalla.core.sdk.model.user.User
import org.akrck02.valhalla.core.sdk.repository.UserRepository
import org.junit.jupiter.api.BeforeEach
import kotlin.test.Test

class UserDataAccessTest : BasicDataAccessTest() {

    companion object {
        lateinit var userRepository: UserRepository
    }

    @BeforeEach
    fun resetDatabaseImpl() {
        userRepository = UserDataAccess(database!!)
    }

    @Test
    fun register() {
        runBlocking {
            userRepository.register(
                User(
                    id = null,
                    username = "akrck01",
                    email = "akrck02@gmail.com",
                    password = "#PasswordisHereLoL#?"
                )
            )
        }
    }

    @Test
    fun get() {
        throw NotImplementedException("not implemented yet!")
    }

    @Test
    fun delete() {
        throw NotImplementedException("not implemented yet!")
    }

    @Test
    fun update() {
        throw NotImplementedException("not implemented yet!")
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
