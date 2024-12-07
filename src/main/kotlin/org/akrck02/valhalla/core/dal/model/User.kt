package org.akrck02.valhalla.core.dal.model

import org.akrck02.valhalla.core.sdk.model.user.User
import org.bson.Document
import org.bson.types.ObjectId

/**
 * Extension function to convert a User to a Document
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
 * Extension function to convert a Document to a User
 * so mongodb can understand.
 * @return Deserialized user
 */
fun Document?.asUser(): User? {
    this ?: return null
    return User(
        id = getObjectId("_id").toHexString(),
        username = getString(User::username.name.lowercase()),
        email = getString(User::email.name.lowercase()),
        password = getString(User::password.name.lowercase()),
        validated = getBoolean(User::validated.name.lowercase()),
        validationCode = getString(User::validationCode.name.lowercase()),
        profilePicturePath = getString(User::profilePicturePath.name.lowercase()),
        creationTime = getLong(User::creationTime.name.lowercase())
    )
}

