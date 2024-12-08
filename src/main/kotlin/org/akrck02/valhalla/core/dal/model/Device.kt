package org.akrck02.valhalla.core.dal.model

import org.akrck02.valhalla.core.sdk.model.device.Device
import org.bson.Document


/**
 * Extension function to convert a device to a document
 * so mongodb can understand.
 * @return Serialized mongodb document
 */
fun Device?.asDocument(): Document? {

    this ?: return null
    val doc = Document()
        .append(Device::address.name.lowercase(), address)
        .append(Device::userAgent.name.lowercase(), userAgent)
        .append(Device::token.name.lowercase(), token)

    return doc
}

/**
 * Extension function to convert a document to a device
 * so mongodb can understand.
 * @return Deserialized device
 */
fun Document?.asDevice(): Device? {
    this ?: return null
    return Device(
        address = getString(Device::address.name.lowercase()),
        userAgent = getString(Device::userAgent.name.lowercase()),
        token = getString(Device::token.name.lowercase())
    )
}

