package org.akrck02.valhalla.core.dal.model

import com.mongodb.client.model.Updates
import org.akrck02.valhalla.core.sdk.model.device.Device
import org.akrck02.valhalla.core.sdk.model.user.User
import org.bson.Document
import org.bson.conversions.Bson
import org.bson.types.ObjectId


/**
 * Extension function to convert a device to a document
 * so mongodb can understand.
 * @return Serialized mongodb document
 */
fun Device?.asDocument(): Document? {

    this ?: return null
    val doc = Document()
        .append(Device::owner.name.lowercase(), owner)
        .append(Device::address.name.lowercase(), address)
        .append(Device::userAgent.name.lowercase(), userAgent)
        .append(Device::token.name.lowercase(), token)

    id?.let {
        doc.append("_id", id ?: ObjectId(id))
    }

    return doc
}

/**
 * Get the updates to be done from one device to match the other
 * @param other The device to compare with
 * @return BSON document with the changes or null
 */
fun Device.getUpdatesToBeDone(other: Device): Bson? {

    val updates = mutableListOf<Bson>()

    if (other.token != this.token) {
        updates.add(Updates.set(Device::token.name.lowercase(), other.token))
    }

    return if (updates.isEmpty()) null else Updates.combine(updates)
}

/**
 * Extension function to convert a document to a device
 * so mongodb can understand.
 * @return Deserialized device
 */
fun Document?.asDevice(): Device? {
    this ?: return null
    return Device(
        id = getObjectId("_id").toHexString(),
        owner = get(Device::owner.name.lowercase()) as User,
        address = getString(Device::address.name.lowercase()),
        userAgent = getString(Device::userAgent.name.lowercase()),
        token = getString(Device::token.name.lowercase())
    )
}

