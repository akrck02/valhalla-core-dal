package org.akrck02.valhalla.core.dal.service.user

import com.mongodb.client.model.Filters
import com.mongodb.kotlin.client.coroutine.MongoDatabase
import kotlinx.coroutines.flow.firstOrNull
import org.akrck02.valhalla.core.dal.database.DatabaseCollections
import org.akrck02.valhalla.core.dal.database.mongoIdEquals
import org.akrck02.valhalla.core.dal.model.asDocument
import org.akrck02.valhalla.core.dal.model.asUser
import org.akrck02.valhalla.core.dal.model.getUpdatesToBeDone
import org.akrck02.valhalla.core.dal.model.validateCompulsoryProperties
import org.akrck02.valhalla.core.sdk.error.ErrorCode
import org.akrck02.valhalla.core.sdk.model.device.Device
import org.akrck02.valhalla.core.sdk.model.exception.ServiceException
import org.akrck02.valhalla.core.sdk.model.http.HttpStatusCode
import org.akrck02.valhalla.core.sdk.model.user.User
import org.akrck02.valhalla.core.sdk.repository.UserRepository
import org.bson.Document
import org.bson.types.ObjectId

/**
 * This class represents the data access layer
 * for the user collection.
 */
class UserDataAccess(private val database: MongoDatabase) : UserRepository {

    // region crud

    override suspend fun register(user: User?): String {

        user.validateCompulsoryProperties()

        val userCollection = database.getCollection<Document>(DatabaseCollections.Users.id)
        val sameMailFilter = Filters.eq(User::email.name, user!!.email)
        val existingUser: User? = userCollection.withDocumentClass<Document>()
            .find(sameMailFilter)
            .firstOrNull()
            .asUser()

        existingUser?.let {
            throw ServiceException(
                status = HttpStatusCode.BadRequest,
                code = ErrorCode.UserAlreadyExists,
                message = "User already exists."
            )
        }

        user.id = null
        val insertedId = userCollection.insertOne(user.asDocument()!!).insertedId
        insertedId ?: throw ServiceException(
            status = HttpStatusCode.InternalServerError,
            code = ErrorCode.DatabaseError,
            message = "User could not be added."
        )

        return insertedId.asObjectId().value.toHexString()
    }


    override suspend fun get(id: String?, secure: Boolean?): User {

        id.takeIf { it.isNullOrBlank().not() } ?: throw ServiceException(
            status = HttpStatusCode.BadRequest,
            code = ErrorCode.InvalidRequest,
            message = "User id cannot be empty."
        )

        val userCollection = database.getCollection<User>(DatabaseCollections.Users.id)
        val user: User? = userCollection.withDocumentClass<Document>()
            .find(mongoIdEquals(id))
            .firstOrNull()
            .asUser()

        user ?: throw ServiceException(
            status = HttpStatusCode.BadRequest,
            code = ErrorCode.NotFound,
            message = "User with id $id does not exist."
        )

        return user.apply { if (secure == true) password = null }
    }

    override suspend fun getByEmail(email: String?, secure: Boolean?): User {

        email.takeIf { it.isNullOrBlank().not() } ?: throw ServiceException(
            status = HttpStatusCode.BadRequest,
            code = ErrorCode.InvalidRequest,
            message = "User email cannot be empty."
        )

        val userCollection = database.getCollection<User>(DatabaseCollections.Users.id)
        val user: User? = userCollection.withDocumentClass<Document>()
            .find(Filters.eq(User::email.name, email))
            .firstOrNull()
            .asUser()

        user ?: throw ServiceException(
            status = HttpStatusCode.BadRequest,
            code = ErrorCode.NotFound,
            message = "User with email $email does not exist."
        )

        return user.apply { if (secure == true) password = null }
    }

    override suspend fun delete(id: String?) {

        id.takeIf { it.isNullOrBlank().not() } ?: throw ServiceException(
            status = HttpStatusCode.BadRequest,
            code = ErrorCode.InvalidRequest,
            message = "User id cannot be empty."
        )

        get(id)

        val userCollection = database.getCollection<User>(DatabaseCollections.Users.id)
        val result = userCollection.deleteOne(Filters.eq("_id", ObjectId(id)))

        takeIf { result.deletedCount > 0 } ?: throw ServiceException(
            status = HttpStatusCode.InternalServerError,
            code = ErrorCode.DatabaseError,
            message = "User could not be deleted."
        )

    }

    override suspend fun update(id: String?, user: User?) {

        id.takeIf { it.isNullOrBlank().not() } ?: throw ServiceException(
            status = HttpStatusCode.BadRequest,
            code = ErrorCode.InvalidRequest,
            message = "User id cannot be empty."
        )

        user.validateCompulsoryProperties()

        val idFilter = Filters.eq("_id", ObjectId(id))
        user!!.id = id

        val userCollection = database.getCollection<Document>(DatabaseCollections.Users.id)
        val existingUser: User = get(id, false)
        userCollection.updateOne(idFilter, existingUser.getUpdatesToBeDone(user))

    }

    override suspend fun updateProfilePicture(id: String?, picture: ByteArray?) {
        TODO("Not yet implemented")
    }

    // endregion
    // region login

    override suspend fun login(user: User?, device: Device?): String {
        TODO("Not yet implemented")
    }

    override suspend fun loginWithAuth(id: String?, token: String?) {
        TODO("Not yet implemented")
    }

    // endregion
    // region validate

    override suspend fun validateAccount(code: String?) {
        TODO("Not yet implemented")
    }

    //endregion
}
