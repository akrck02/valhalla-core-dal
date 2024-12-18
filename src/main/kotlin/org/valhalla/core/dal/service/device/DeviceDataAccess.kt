package org.valhalla.core.dal.service.device

import com.mongodb.client.model.Filters.and
import com.mongodb.client.model.Filters.eq
import com.mongodb.client.model.Projections
import com.mongodb.kotlin.client.coroutine.MongoDatabase
import kotlinx.coroutines.flow.firstOrNull
import org.valhalla.core.dal.database.DatabaseCollections
import org.valhalla.core.dal.database.idEqualsFilter
import org.valhalla.core.sdk.error.ErrorCode
import org.valhalla.core.sdk.model.device.Device
import org.valhalla.core.sdk.model.exception.ServiceException
import org.valhalla.core.sdk.model.http.HttpStatusCode
import org.valhalla.core.sdk.model.user.User
import org.valhalla.core.sdk.repository.DeviceRepository
import org.valhalla.core.sdk.repository.UserRepository

class DeviceDataAccess(
    private val database: MongoDatabase,
    private val userRepository: UserRepository
) : DeviceRepository {

    override suspend fun register(userId: String?, device: Device?): String {

        userId ?: throw ServiceException(
            status = HttpStatusCode.BadRequest,
            code = ErrorCode.InvalidRequest,
            message = "User id cannot be empty."
        )

        device ?: throw ServiceException(
            status = HttpStatusCode.BadRequest,
            code = ErrorCode.InvalidRequest,
            message = "Device cannot be empty."
        )

        device.token ?: throw ServiceException(
            status = HttpStatusCode.BadRequest,
            code = ErrorCode.InvalidRequest,
            message = "Device token cannot be empty."
        )

        var user = userRepository.get(userId)
        var foundDevice = user.devices.find { (it.userAgent == device.userAgent).and(it.address == device.address) }

        if (null == foundDevice) {
            user.devices.add(device)
        } else {
            foundDevice.apply { token = device.token }
        }

        userRepository.update(userId, user)
        user = userRepository.get(userId)

        foundDevice = user.devices.find { (it.userAgent == device.userAgent).and(it.address == device.address) }
        foundDevice ?: throw ServiceException(
            status = HttpStatusCode.BadRequest,
            code = ErrorCode.InvalidRequest,
            message = "Device not found."
        )

        return foundDevice.token!!
    }

    override suspend fun get(userId: String?, id: String?): Device {

        userId ?: throw ServiceException(
            status = HttpStatusCode.BadRequest,
            code = ErrorCode.InvalidRequest,
            message = "User id cannot be empty."
        )

        id ?: throw ServiceException(
            status = HttpStatusCode.BadRequest,
            code = ErrorCode.InvalidRequest,
            message = "Device id cannot be empty."
        )

        TODO("Not yet implemented")
    }

    override suspend fun getAll(userId: String?): List<Device> {

        userId ?: throw ServiceException(
            status = HttpStatusCode.BadRequest,
            code = ErrorCode.InvalidRequest,
            message = "User id cannot be empty."
        )

        TODO("Not yet implemented")
    }

    override suspend fun getByAuth(userId: String?, token: String?): Device {

        userId ?: throw ServiceException(
            status = HttpStatusCode.BadRequest,
            code = ErrorCode.InvalidRequest,
            message = "User id cannot be empty."
        )

        token ?: throw ServiceException(
            status = HttpStatusCode.BadRequest,
            code = ErrorCode.InvalidRequest,
            message = "Token cannot be empty."
        )

        val users = database.getCollection<User>(DatabaseCollections.Users.id)
        val foundUserDevice = users.find<User>(filter = and(idEqualsFilter(userId), eq("devices.token", token)))
            .projection(
                Projections.include(
                    "devices.token",
                    "devices.userAgent",
                    "devices.address"
                )
            )
            .firstOrNull()

        return foundUserDevice?.devices?.firstOrNull() ?: throw ServiceException(
            status = HttpStatusCode.BadRequest,
            code = ErrorCode.NotFound,
            message = "Device not found."
        )
    }
}