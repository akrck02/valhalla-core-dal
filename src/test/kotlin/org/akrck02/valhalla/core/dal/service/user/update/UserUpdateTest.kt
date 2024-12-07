package org.akrck02.valhalla.core.dal.service.user.update

import kotlinx.coroutines.runBlocking
import org.akrck02.valhalla.core.dal.mock.CorrectUser
import org.akrck02.valhalla.core.dal.service.user.UserDataAccess
import org.akrck02.valhalla.core.dal.tool.BaseDataAccessTest
import org.akrck02.valhalla.core.sdk.repository.UserRepository
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.Test
import kotlin.test.assertEquals

class UserUpdateTest : BaseDataAccessTest() {

    companion object {
        lateinit var userRepository: UserRepository
    }

    @BeforeEach
    fun resetDatabaseImpl() {
        userRepository = UserDataAccess(database!!)
    }

    @Test
    fun `update (happy path)`() = runBlocking {

        val originalUser = CorrectUser.copy()
        originalUser.id = userRepository.register(originalUser)

        val updatedUser = CorrectUser.copy(username = "xxx_shadow_the_hedgehog_xxx")
        userRepository.update(originalUser.id, updatedUser)

        val databaseUser = userRepository.get(originalUser.id, false)
        assertEquals(updatedUser, databaseUser)

    }
}