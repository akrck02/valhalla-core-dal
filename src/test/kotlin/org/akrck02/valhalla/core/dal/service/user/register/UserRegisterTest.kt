package org.akrck02.valhalla.core.dal.service.user.register

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

class UserRegisterTest : BaseDataAccessTest() {

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
            val id = userRepository.register(CorrectUser)
            println("Inserted user with id $id")
        }
    }

    @Test
    fun registerEmptyUsername() {
        assertThrowsServiceException(
            ServiceException(
                status = HttpStatusCode.BadRequest,
                code = ErrorCode.InvalidRequest,
                message = "Username cannot be empty.",
            )
        ) {
            runBlocking {
                userRepository.register(
                    CorrectUser.copy(username = null)
                )
            }
        }
    }

    @Test
    fun registerEmptyEmail() {
        assertThrowsServiceException(
            ServiceException(
                status = HttpStatusCode.BadRequest,
                code = ErrorCode.InvalidRequest,
                message = "Email cannot be empty.",
            )
        ) {
            runBlocking {
                userRepository.register(
                    CorrectUser.copy(email = null)
                )
            }
        }
    }

    @Test
    fun registerEmptyPassword() {
        assertThrowsServiceException(
            ServiceException(
                status = HttpStatusCode.BadRequest,
                code = ErrorCode.InvalidRequest,
                message = "Password cannot be empty.",
            )
        ) {
            runBlocking {
                userRepository.register(
                    CorrectUser.copy(password = null)
                )
            }
        }
    }
}