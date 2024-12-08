package org.akrck02.valhalla.core.dal.model

import com.mongodb.client.model.Updates
import org.akrck02.valhalla.core.sdk.error.ErrorCode
import org.akrck02.valhalla.core.sdk.model.exception.ServiceException
import org.akrck02.valhalla.core.sdk.model.http.HttpStatusCode
import org.akrck02.valhalla.core.sdk.model.user.User
import org.akrck02.valhalla.core.sdk.validation.validateEmail
import org.akrck02.valhalla.core.sdk.validation.validatePassword
import org.bson.Document
import org.bson.conversions.Bson
import org.bson.types.ObjectId

/**
 * Extension function to validate compulsory
 * properties for a user
 * @param validatePassword If password was validated
 * @throws ServiceException if a requirement is not being fulfilled
 */
fun User?.validateCompulsoryProperties(validatePassword: Boolean = true) {

    this ?: throw ServiceException(
        status = HttpStatusCode.BadRequest,
        code = ErrorCode.InvalidRequest,
        message = "User cannot be empty."
    )

    takeIf { it.email.isNullOrBlank().not() } ?: throw ServiceException(
        status = HttpStatusCode.BadRequest,
        code = ErrorCode.InvalidEmail,
        message = "Email cannot be empty."
    )

    if (validatePassword) {
        takeIf { it.password.isNullOrBlank().not() } ?: throw ServiceException(
            status = HttpStatusCode.BadRequest,
            code = ErrorCode.InvalidPassword,
            message = "Password cannot be empty."
        )
    }

    takeIf { it.username.isNullOrBlank().not() } ?: throw ServiceException(
        status = HttpStatusCode.BadRequest,
        code = ErrorCode.InvalidRequest,
        message = "Username cannot be empty."
    )

    email?.validateEmail()
    if (validatePassword) {
        password?.validatePassword()
    }
}

/**
 * Extension function to convert a user to a document
 * so mongodb can understand.
 * @return Serialized mongodb document
 */
fun User?.asDocument(): Document? {

    this ?: return null
    val doc = Document()
        .append(User::username.name.lowercase(), username)
        .append(User::email.name.lowercase(), email)
        .append(User::password.name.lowercase(), password)
        .append(User::validated.name.lowercase(), validated)
        .append(User::validationCode.name.lowercase(), validationCode)
        .append(User::profilePicturePath.name.lowercase(), profilePicturePath)
        .append(User::creationTime.name.lowercase(), creationTime)

    id?.let {
        doc.append("_id", id ?: ObjectId(id))
    }
    return doc
}

/**
 * Get the updates to be done from one user to match the other
 * @param other The user to compare with
 * @return BSON document with the changes or null
 */
fun User.getUpdatesToBeDone(other: User): Bson? {

    val updates = mutableListOf<Bson>()

    if (other.username != this.username) {
        updates.add(Updates.set(User::username.name.lowercase(), other.username))
    }

    var devicesChanged = false
    other.devices.forEach {
        if (this.devices.contains(it).not()) {
            devicesChanged = true
        }
    }

    if (devicesChanged) {
        updates.add(Updates.set(User::devices.name.lowercase(), other.devices))
    }

    return if (updates.isEmpty()) null else Updates.combine(updates)
}

/**
 * Extension function to convert a document to a user
 * so mongodb can understand.
 * @return Deserialized user
 */
fun Document?.asUser(): User? {
    this ?: return null
    val user = User(
        id = getObjectId("_id").toHexString(),
        username = getString(User::username.name.lowercase()),
        email = getString(User::email.name.lowercase()),
        password = getString(User::password.name.lowercase()),
        validated = getBoolean(User::validated.name.lowercase()),
        validationCode = getString(User::validationCode.name.lowercase()),
        profilePicturePath = getString(User::profilePicturePath.name.lowercase()),
        creationTime = getLong(User::creationTime.name.lowercase())
    )


    getList(User::devices.name.lowercase(), Document::class.java)?.forEach {
        user.devices.add(it.asDevice()!!)
    }

    return user
}



