package org.akrck02.valhalla.core.dal.service.user

import com.mongodb.client.model.Filters
import com.mongodb.kotlin.client.coroutine.MongoDatabase
import kotlinx.coroutines.flow.firstOrNull
import org.akrck02.valhalla.core.dal.database.DatabaseCollections
import org.akrck02.valhalla.core.sdk.error.ErrorCode
import org.akrck02.valhalla.core.sdk.model.device.Device
import org.akrck02.valhalla.core.sdk.model.exception.ServiceException
import org.akrck02.valhalla.core.sdk.model.http.HttpStatusCode
import org.akrck02.valhalla.core.sdk.model.user.User
import org.akrck02.valhalla.core.sdk.repository.UserRepository
import org.akrck02.valhalla.core.sdk.validation.validateEmail
import org.akrck02.valhalla.core.sdk.validation.validatePassword

/**
 *
 */
class UserDataAccess(val database: MongoDatabase) : UserRepository {

    override suspend fun delete(id: String?) {
        TODO("Not yet implemented")
    }

    override suspend fun get(id: String?): User? {
        TODO("Not yet implemented")
    }

    override suspend fun getByEmail(email: String?): User? {
        TODO("Not yet implemented")
    }

    override suspend fun login(user: User?, device: Device?): String? {
        TODO("Not yet implemented")
    }

    override suspend fun loginWithAuth(id: String?, token: String?) {
        TODO("Not yet implemented")
    }

    override suspend fun register(user: User?) {
        user ?: throw ServiceException(
            status = HttpStatusCode.BadRequest,
            code = ErrorCode.InvalidRequest,
            message = "User cannot be empty."
        )

        user.takeIf { (it.email ?: "").isNotBlank() } ?: throw ServiceException(
            status = HttpStatusCode.BadRequest,
            code = ErrorCode.InvalidEmail,
            message = "Email cannot be empty."
        )

        user.takeIf { (it.password ?: "").isNotBlank() } ?: throw ServiceException(
            status = HttpStatusCode.BadRequest,
            code = ErrorCode.InvalidPassword,
            message = "Password cannot be empty."
        )

        user.takeIf { (it.username ?: "").isNotBlank() } ?: throw ServiceException(
            status = HttpStatusCode.BadRequest,
            code = ErrorCode.InvalidRequest,
            message = "Username cannot be empty."
        )

        user.email?.validateEmail()
        user.password?.validatePassword()

        val userCollection = database.getCollection<User>(DatabaseCollections.Users.id)
        val existingUser: User? = userCollection.withDocumentClass<User>()
            .find(Filters.eq(User::username.name, user.username))
            .firstOrNull()

        existingUser?.let {
            throw ServiceException(
                status = HttpStatusCode.BadRequest,
                code = ErrorCode.UserAlreadyExists,
                message = "User already exists."
            )
        }

        userCollection.insertOne(user)
    }

    override suspend fun update(id: String?, user: User?) {
        TODO("Not yet implemented")
    }

    override suspend fun updateProfilePicture(id: String?, picture: ByteArray?) {
        TODO("Not yet implemented")
    }

    override suspend fun validateAccount(code: String?) {
        TODO("Not yet implemented")
    }

}
