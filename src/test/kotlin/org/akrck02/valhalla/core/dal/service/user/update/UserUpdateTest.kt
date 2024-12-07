package org.akrck02.valhalla.core.dal.service.user.update

import kotlinx.coroutines.runBlocking
import org.akrck02.valhalla.core.dal.mock.CorrectUser
import org.akrck02.valhalla.core.dal.service.user.UserDataAccess
import org.akrck02.valhalla.core.dal.tool.BaseDataAccessTest
import org.akrck02.valhalla.core.dal.tool.assertThrowsServiceException
import org.akrck02.valhalla.core.sdk.error.ErrorCode
import org.akrck02.valhalla.core.sdk.model.exception.ServiceException
import org.akrck02.valhalla.core.sdk.model.http.HttpStatusCode
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

    @Test
    fun `update nothing changed (happy path)`() = runBlocking {

        val originalUser = CorrectUser.copy()
        originalUser.id = userRepository.register(originalUser)

        assertThrowsServiceException(
            ServiceException(
                status = HttpStatusCode.BadRequest,
                code = ErrorCode.NothingChanged,
                message = "No changes needed.",
            )
        ) {
            userRepository.update(originalUser.id, originalUser)
        }
    }

    @Test
    fun `update with empty user`(): Unit = runBlocking {

        val originalUser = CorrectUser.copy()
        originalUser.id = userRepository.register(originalUser)

        assertThrowsServiceException(
            ServiceException(
                status = HttpStatusCode.BadRequest,
                code = ErrorCode.NothingChanged,
                message = "User cannot be empty.",
            )
        ) {
            userRepository.update(originalUser.id, null)
        }
    }


    @Test
    fun `update with empty email`(): Unit = runBlocking {

        val originalUser = CorrectUser.copy()
        originalUser.id = userRepository.register(originalUser)

        assertThrowsServiceException(
            ServiceException(
                status = HttpStatusCode.BadRequest,
                code = ErrorCode.NothingChanged,
                message = "Email cannot be empty.",
            )
        ) {
            userRepository.update(originalUser.id, originalUser.copy(email = ""))
        }

        assertThrowsServiceException(
            ServiceException(
                status = HttpStatusCode.BadRequest,
                code = ErrorCode.NothingChanged,
                message = "Email cannot be empty.",
            )
        ) {
            userRepository.update(originalUser.id, originalUser.copy(email = null))
        }
    }

    @Test
    fun `update with empty username`(): Unit = runBlocking {

        val originalUser = CorrectUser.copy()
        originalUser.id = userRepository.register(originalUser)

        assertThrowsServiceException(
            ServiceException(
                status = HttpStatusCode.BadRequest,
                code = ErrorCode.NothingChanged,
                message = "Username cannot be empty.",
            )
        ) {
            userRepository.update(originalUser.id, originalUser.copy(username = ""))
        }

        assertThrowsServiceException(
            ServiceException(
                status = HttpStatusCode.BadRequest,
                code = ErrorCode.NothingChanged,
                message = "Username cannot be empty.",
            )
        ) {
            userRepository.update(originalUser.id, originalUser.copy(username = null))
        }
    }

}