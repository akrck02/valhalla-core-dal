package org.akrck02.valhalla.core.dal.service.user.get

import kotlinx.coroutines.runBlocking
import org.akrck02.valhalla.core.dal.mock.CorrectUser
import org.akrck02.valhalla.core.dal.service.user.UserDataAccess
import org.akrck02.valhalla.core.dal.tool.BaseDataAccessTest
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
    fun `get user (happy path)`() = runBlocking {

        val user = CorrectUser.copy()
        user.id = userRepository.register(user)
        println("Inserted user with id ${user.id}.")

        val foundUser = userRepository.get(user.id, secure = false)
        assertEquals(user, foundUser)
        println("User ${user.username} (id: ${user.id}) found.")

    }

    @Test
    fun `get user that does not exists`() = runBlocking {


    }
}