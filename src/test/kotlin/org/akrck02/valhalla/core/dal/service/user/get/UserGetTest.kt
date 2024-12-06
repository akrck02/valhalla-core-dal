package org.akrck02.valhalla.core.dal.service.user.get

import kotlinx.coroutines.runBlocking
import org.akrck02.valhalla.core.dal.service.user.UserDataAccess
import org.akrck02.valhalla.core.dal.tool.BaseDataAccessTest
import org.akrck02.valhalla.core.sdk.model.user.User
import org.akrck02.valhalla.core.sdk.repository.UserRepository
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.Test
import kotlin.test.assertEquals

class UserGetTest : BaseDataAccessTest() {

    companion object {
        lateinit var userRepository: UserRepository
    }

    @BeforeEach
    fun resetDatabaseImpl() {
        userRepository = UserDataAccess(database!!)
    }

    @Test
    fun get() {
        runBlocking {

            val user = User(
                id = null,
                username = "akrck01",
                email = "akrck02@gmail.com",
                password = "#PasswordisHereLoL#?"
            )

            user.id = userRepository.register(user)
            val foundUser = userRepository.get(user.id, secure = false)

            assertEquals(user, foundUser)

        }
    }

}