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
import org.akrck02.valhalla.core.sdk.validation.MINIMUM_CHARACTERS_FOR_EMAIL
import org.akrck02.valhalla.core.sdk.validation.MINIMUM_CHARACTERS_FOR_PASSWORD
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
    fun registerNoAtEmail() {
        assertThrowsServiceException(
            ServiceException(
                status = HttpStatusCode.BadRequest,
                code = ErrorCode.InvalidRequest,
                message = "Email must have one @.",
            )
        ) {
            runBlocking {
                userRepository.register(
                    CorrectUser.copy(email = "akrck02.com")
                )
            }
        }
    }

    @Test
    fun registerNoDotEmail() {
        assertThrowsServiceException(
            ServiceException(
                status = HttpStatusCode.BadRequest,
                code = ErrorCode.InvalidRequest,
                message = "Email must have at least one dot.",
            )
        ) {
            runBlocking {
                userRepository.register(
                    CorrectUser.copy(email = "akrck02@com")
                )
            }
        }
    }

    @Test
    fun registerShortEmail() {
        assertThrowsServiceException(
            ServiceException(
                status = HttpStatusCode.BadRequest,
                code = ErrorCode.InvalidRequest,
                message = "Email must have at least $MINIMUM_CHARACTERS_FOR_EMAIL characters.",
            )
        ) {
            runBlocking {
                userRepository.register(
                    CorrectUser.copy(email = "com")
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

    @Test
    fun registerShortPassword() {
        assertThrowsServiceException(
            ServiceException(
                status = HttpStatusCode.BadRequest,
                code = ErrorCode.InvalidRequest,
                message = "Password must have at least $MINIMUM_CHARACTERS_FOR_PASSWORD characters.",
            )
        ) {
            runBlocking {
                userRepository.register(
                    CorrectUser.copy(password = "123")
                )
            }
        }
    }

    @Test
    fun registerNotAlphanumericPassword() {
        assertThrowsServiceException(
            ServiceException(
                status = HttpStatusCode.BadRequest,
                code = ErrorCode.InvalidRequest,
                message = "Password must have at least $MINIMUM_CHARACTERS_FOR_PASSWORD characters.",
            )
        ) {
            runBlocking {
                userRepository.register(
                    CorrectUser.copy(password = "123456789AbCdEf*")
                )
            }
        }
    }

    @Test
    fun registerNotUppercasePassword() {
        assertThrowsServiceException(
            ServiceException(
                status = HttpStatusCode.BadRequest,
                code = ErrorCode.InvalidRequest,
                message = "Password must have at least one uppercase character.",
            )
        ) {
            runBlocking {
                userRepository.register(
                    CorrectUser.copy(password = "#aeeeeeeeeeeeeeee")
                )
            }
        }
    }

    @Test
    fun registerNotLowercasePassword() {
        assertThrowsServiceException(
            ServiceException(
                status = HttpStatusCode.BadRequest,
                code = ErrorCode.InvalidRequest,
                message = "Password must have at least one lowercase character.",
            )
        ) {
            runBlocking {
                userRepository.register(
                    CorrectUser.copy(password = "#AEEEEEEEEEEEEEE")
                )
            }
        }
    }
}