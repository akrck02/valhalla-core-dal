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
    fun `register (happy path)`() = runBlocking {
        runBlocking {
            val id = userRepository.register(CorrectUser)
            println("Inserted user with id $id")
        }
    }

    @Test
    fun `register with empty username`() = runBlocking {

        assertThrowsServiceException(
            ServiceException(
                status = HttpStatusCode.BadRequest,
                code = ErrorCode.InvalidRequest,
                message = "Username cannot be empty.",
            )
        ) {
            userRepository.register(CorrectUser.copy(username = null))
        }

        assertThrowsServiceException(
            ServiceException(
                status = HttpStatusCode.BadRequest,
                code = ErrorCode.InvalidRequest,
                message = "Username cannot be empty.",
            )
        ) {
            userRepository.register(CorrectUser.copy(username = ""))
        }
    }

    @Test
    fun `register with empty email`() = runBlocking {
        assertThrowsServiceException(
            ServiceException(
                status = HttpStatusCode.BadRequest,
                code = ErrorCode.InvalidRequest,
                message = "Email cannot be empty.",
            )
        ) {
            userRepository.register(CorrectUser.copy(email = null))
        }

        assertThrowsServiceException(
            ServiceException(
                status = HttpStatusCode.BadRequest,
                code = ErrorCode.InvalidRequest,
                message = "Email cannot be empty.",
            )
        ) {
            userRepository.register(CorrectUser.copy(email = ""))
        }
    }

    @Test
    fun `register with not at email`() = runBlocking {
        assertThrowsServiceException(
            ServiceException(
                status = HttpStatusCode.BadRequest,
                code = ErrorCode.InvalidRequest,
                message = "Email must have one @.",
            )
        ) {
            userRepository.register(CorrectUser.copy(email = "akrck02.com"))
        }
    }

    @Test
    fun `register with not dot email`() = runBlocking {
        assertThrowsServiceException(
            ServiceException(
                status = HttpStatusCode.BadRequest,
                code = ErrorCode.InvalidRequest,
                message = "Email must have at least one dot.",
            )
        ) {
            userRepository.register(CorrectUser.copy(email = "akrck02@com"))
        }
    }

    @Test
    fun `register with short email`() = runBlocking {
        assertThrowsServiceException(
            ServiceException(
                status = HttpStatusCode.BadRequest,
                code = ErrorCode.InvalidRequest,
                message = "Email must have at least $MINIMUM_CHARACTERS_FOR_EMAIL characters.",
            )
        ) {
            userRepository.register(CorrectUser.copy(email = "com"))
        }
    }

    @Test
    fun `register with empty password`() = runBlocking {
        assertThrowsServiceException(
            ServiceException(
                status = HttpStatusCode.BadRequest,
                code = ErrorCode.InvalidRequest,
                message = "Password cannot be empty.",
            )
        ) {
            userRepository.register(CorrectUser.copy(password = null))
        }

        assertThrowsServiceException(
            ServiceException(
                status = HttpStatusCode.BadRequest,
                code = ErrorCode.InvalidRequest,
                message = "Password cannot be empty.",
            )
        ) {
            userRepository.register(CorrectUser.copy(password = ""))
        }
    }

    @Test
    fun `register with short password`() = runBlocking {
        assertThrowsServiceException(
            ServiceException(
                status = HttpStatusCode.BadRequest,
                code = ErrorCode.InvalidRequest,
                message = "Password must have at least $MINIMUM_CHARACTERS_FOR_PASSWORD characters.",
            )
        ) {
            userRepository.register(CorrectUser.copy(password = "123"))
        }
    }

    @Test
    fun `register with not alphanumeric password`() = runBlocking {
        assertThrowsServiceException(
            ServiceException(
                status = HttpStatusCode.BadRequest,
                code = ErrorCode.InvalidRequest,
                message = "Password must have at least one number.",
            )
        ) {
            userRepository.register(CorrectUser.copy(password = "#Aaaaaaaaaaaaaaa"))
        }
    }

    @Test
    fun `register with not uppercase password`() = runBlocking {
        assertThrowsServiceException(
            ServiceException(
                status = HttpStatusCode.BadRequest,
                code = ErrorCode.InvalidRequest,
                message = "Password must have at least one lowercase and uppercase character.",
            )
        ) {
            userRepository.register(CorrectUser.copy(password = "aeeeeeeeeeeeeeee"))
        }
    }

    @Test
    fun `register with not lowercase password`() = runBlocking {
        assertThrowsServiceException(
            ServiceException(
                status = HttpStatusCode.BadRequest,
                code = ErrorCode.InvalidRequest,
                message = "Password must have at least one lowercase and uppercase character.",
            )
        ) {
            userRepository.register(CorrectUser.copy(password = "AEEEEEEEEEEEEEEE"))
        }
    }

    @Test
    fun `register with not special character password`(): Unit = runBlocking {
        assertThrowsServiceException(
            ServiceException(
                status = HttpStatusCode.BadRequest,
                code = ErrorCode.InvalidRequest,
                message = "Password must have at least one special character.",
            )
        ) {
            userRepository.register(CorrectUser.copy(password = "AmazingPassword2000"))
        }

        userRepository.register(CorrectUser.copy(email = "0001@mail.com", password = "#AmazingPassword2000"))
        userRepository.register(CorrectUser.copy(email = "0002@mail.com", password = "*AmazingPassword2000"))
        userRepository.register(CorrectUser.copy(email = "0004@mail.com", password = "?AmazingPassword2000"))
        userRepository.register(CorrectUser.copy(email = "0005@mail.com", password = "¿AmazingPassword2000"))
        userRepository.register(CorrectUser.copy(email = "0006@mail.com", password = "¡AmazingPassword2000"))
        userRepository.register(CorrectUser.copy(email = "0007@mail.com", password = "!AmazingPassword2000"))
        userRepository.register(CorrectUser.copy(email = "0008@mail.com", password = "&AmazingPassword2000"))
        userRepository.register(CorrectUser.copy(email = "0009@mail.com", password = "^AmazingPassword2000"))
        userRepository.register(CorrectUser.copy(email = "0010@mail.com", password = "\$AmazingPassword2000"))
        userRepository.register(CorrectUser.copy(email = "0011@mail.com", password = "%AmazingPassword2000"))

    }

    @Test
    fun `register user that already exists`() = runBlocking {

        userRepository.register(CorrectUser)
        assertThrowsServiceException(
            ServiceException(
                status = HttpStatusCode.BadRequest,
                code = ErrorCode.UserAlreadyExists,
                message = "User already exists.",
            )
        ) {
            userRepository.register(CorrectUser)
        }

    }
}