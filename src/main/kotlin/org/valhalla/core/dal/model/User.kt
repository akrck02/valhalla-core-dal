package org.valhalla.core.dal.model

import com.mongodb.client.model.Updates

import org.bson.conversions.Bson
import org.valhalla.core.sdk.error.ErrorCode
import org.valhalla.core.sdk.model.exception.ServiceException
import org.valhalla.core.sdk.model.http.HttpStatusCode
import org.valhalla.core.sdk.model.user.User
import org.valhalla.core.sdk.validation.validateEmail
import org.valhalla.core.sdk.validation.validatePassword
import java.security.MessageDigest

private const val DEFAULT_ENCRYPTION_METHOD = "SHA-256"

/**
 * Extension function to validate compulsory
 * properties for a user.
 * Throws [ServiceException] if a requirement is not being fulfilled.
 */
@Suppress("unused")
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


/** Get the updates to be done from one user to match the other [User] data in [Bson] format. */
@Suppress("unused")
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

/** Hash the [User] password */
@Suppress("unused")
fun User.hashPassword() {

    password ?: return

    val sha256 = MessageDigest.getInstance(DEFAULT_ENCRYPTION_METHOD)
    sha256.update(password!!.toByteArray())
    password = sha256.digest().joinToString("")
}