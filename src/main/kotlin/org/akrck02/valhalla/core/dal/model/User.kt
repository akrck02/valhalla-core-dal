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
        .append(User::username.name, username)
        .append(User::email.name, email)
        .append(User::password.name, password)
        .append(User::validated.name, validated)
        .append(User::validationCode.name, validationCode)
        .append(User::profilePicturePath.name, profilePicturePath)
        .append(User::creationTime.name, creationTime)

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
        username = getString(User::username.name),
        email = getString(User::email.name),
        password = getString(User::password.name),
        validated = getBoolean(User::validated.name),
        validationCode = getString(User::validationCode.name),
        profilePicturePath = getString(User::profilePicturePath.name),
        creationTime = getLong(User::creationTime.name)
    )
}

