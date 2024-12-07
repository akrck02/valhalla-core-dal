package org.akrck02.valhalla.core.dal.service.user.delete

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

class UserDeleteTest : BaseDataAccessTest() {

    companion object {
        lateinit var userRepository: UserRepository
    }

    @BeforeEach
    fun resetDatabaseImpl() {
        userRepository = UserDataAccess(database!!)
    }

    @Test
    fun `delete by id (happy path)`() = runBlocking {

        val user = CorrectUser.copy()
        user.id = userRepository.register(user)
        userRepository.delete(user.id)

        assertThrowsServiceException(
            ServiceException(
                status = HttpStatusCode.BadRequest,
                code = ErrorCode.InvalidRequest,
                message = "User with id ${user.id} does not exist.",
            )
        ) {
            userRepository.get(user.id)
        }

    }

    @Test
    fun `delete by empty id`() = runBlocking {

        assertThrowsServiceException(
            ServiceException(
                status = HttpStatusCode.BadRequest,
                code = ErrorCode.InvalidRequest,
                message = "User id cannot be empty.",
            )
        ) {
            userRepository.delete("")
        }

        assertThrowsServiceException(
            ServiceException(
                status = HttpStatusCode.BadRequest,
                code = ErrorCode.InvalidRequest,
                message = "User id cannot be empty.",
            )
        ) {
            userRepository.delete(null)
        }
    }


    @Test
    fun `delete by not existing id`() = runBlocking {

        assertThrowsServiceException(
            ServiceException(
                status = HttpStatusCode.BadRequest,
                code = ErrorCode.InvalidRequest,
                message = "User with id 6754afd48e7d8e5df2e2105a does not exist.",
            )
        ) {
            userRepository.delete("6754afd48e7d8e5df2e2105a")
        }
    }

}