package org.valhalla.core.dal.service.user.get

import kotlinx.coroutines.runBlocking
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.Test
import org.valhalla.core.dal.mock.TestUser
import org.valhalla.core.dal.service.user.UserDataAccess
import org.valhalla.core.dal.tool.BaseDataAccessTest
import org.valhalla.core.dal.tool.assertThrowsServiceException
import org.valhalla.core.sdk.error.ErrorCode
import org.valhalla.core.sdk.model.exception.ServiceException
import org.valhalla.core.sdk.model.http.HttpStatusCode
import org.valhalla.core.sdk.repository.UserRepository
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
    fun `get user by id (happy path)`() = runBlocking {

        val user = TestUser.copy()
        user.id = userRepository.register(user)
        println("Inserted user with id ${user.id}.")

        val foundUser = userRepository.get(user.id, secure = false)
        assertEquals(user, foundUser)
        println("User ${foundUser.username} (id: ${foundUser.id}) found.")

    }

    @Test
    fun `get secure user by id (happy path)`() = runBlocking {

        val user = TestUser.copy()
        user.id = userRepository.register(user)
        println("Inserted user with id ${user.id}.")

        val foundUser = userRepository.get(user.id, secure = true)
        assertEquals(user.apply { password = null }, foundUser)
        println("User ${foundUser.username} (id: ${foundUser.id}) found.")

    }

    @Test
    fun `get user by empty id`() = runBlocking {

        assertThrowsServiceException(
            ServiceException(
                status = HttpStatusCode.BadRequest,
                code = ErrorCode.InvalidRequest,
                message = "User id cannot be empty.",
            )
        ) {
            userRepository.get("")
        }

        assertThrowsServiceException(
            ServiceException(
                status = HttpStatusCode.BadRequest,
                code = ErrorCode.InvalidRequest,
                message = "User id cannot be empty.",
            )
        ) {
            userRepository.get(null)
        }

    }

    @Test
    fun `get user by not existing id`() = runBlocking {

        val id = "6754a58aeecc751512e130a2"
        assertThrowsServiceException(
            ServiceException(
                status = HttpStatusCode.BadRequest,
                code = ErrorCode.NotFound,
                message = "User with id $id does not exist.",
            )
        ) {
            userRepository.get(id)
        }

    }


    @Test
    fun `get user by email (happy path)`() = runBlocking {

        val user = TestUser.copy()
        user.id = userRepository.register(user)
        println("Inserted user with email ${user.email}.")

        val foundUser = userRepository.getByEmail(user.email, secure = false)
        assertEquals(user, foundUser)
        println("User ${foundUser.username} (email: ${foundUser.email}) found.")

    }

    @Test
    fun `get secure user by email (happy path)`() = runBlocking {

        val user = TestUser.copy()
        user.id = userRepository.register(user)
        println("Inserted user with email ${user.email}.")

        val foundUser = userRepository.getByEmail(user.email, secure = true)
        assertEquals(user.apply { password = null }, foundUser)
        println("User ${foundUser.username} (email: ${foundUser.email}) found.")

    }

    @Test
    fun `get user by empty email`() = runBlocking {

        assertThrowsServiceException(
            ServiceException(
                status = HttpStatusCode.BadRequest,
                code = ErrorCode.InvalidRequest,
                message = "User email cannot be empty.",
            )
        ) {
            userRepository.getByEmail("")
        }

        assertThrowsServiceException(
            ServiceException(
                status = HttpStatusCode.BadRequest,
                code = ErrorCode.InvalidRequest,
                message = "User email cannot be empty.",
            )
        ) {
            userRepository.getByEmail(null)
        }

    }

    @Test
    fun `get user by not existing email`() = runBlocking {

        val email = "thismailisfake@fakeorganitation.org"
        assertThrowsServiceException(
            ServiceException(
                status = HttpStatusCode.BadRequest,
                code = ErrorCode.NotFound,
                message = "User with email $email does not exist.",
            )
        ) {
            userRepository.getByEmail(email)
        }

    }
}