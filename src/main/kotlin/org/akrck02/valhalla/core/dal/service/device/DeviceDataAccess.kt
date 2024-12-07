package org.akrck02.valhalla.core.dal.service.device

import com.mongodb.kotlin.client.coroutine.MongoDatabase
import org.akrck02.valhalla.core.dal.database.DatabaseCollections
import org.akrck02.valhalla.core.sdk.error.ErrorCode
import org.akrck02.valhalla.core.sdk.model.device.Device
import org.akrck02.valhalla.core.sdk.model.exception.ServiceException
import org.akrck02.valhalla.core.sdk.model.http.HttpStatusCode
import org.akrck02.valhalla.core.sdk.repository.DeviceRepository
import org.bson.Document

class DeviceDataAccess(private val database: MongoDatabase) : DeviceRepository {

    override suspend fun register(device: Device?): String {

        device ?: throw ServiceException(
            status = HttpStatusCode.BadRequest,
            code = ErrorCode.InvalidRequest,
            message = "Device cannot be empty."
        )

        val devices = database.getCollection<Document>(DatabaseCollections.Devices.id)

        // TODO:
        //      1. if a device exists with the same address and userAgent for a user
        //          update the auth token and return it.
        //      2. Else create it with a generated token

        return ""
    }

    override suspend fun get(id: String?): Device {
        TODO("Not yet implemented")
    }

    override suspend fun getAllByOwner(userId: String?): List<Device> {
        TODO("Not yet implemented")
    }

    override suspend fun getByAuth(token: String?): Device {
        TODO("Not yet implemented")
    }

    override suspend fun update(id: String?, device: Device?) {
        TODO("Not yet implemented")
    }

}