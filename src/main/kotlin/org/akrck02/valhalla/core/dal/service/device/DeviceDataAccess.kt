package org.akrck02.valhalla.core.dal.service.device

import com.mongodb.kotlin.client.coroutine.MongoDatabase
import org.akrck02.valhalla.core.sdk.error.ErrorCode
import org.akrck02.valhalla.core.sdk.model.device.Device
import org.akrck02.valhalla.core.sdk.model.exception.ServiceException
import org.akrck02.valhalla.core.sdk.model.http.HttpStatusCode
import org.akrck02.valhalla.core.sdk.repository.DeviceRepository
import org.akrck02.valhalla.core.sdk.repository.UserRepository

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

        var user = userRepository.get(userId)
        var foundDevice = user.devices.find { (it.userAgent == device.userAgent).and(it.address == device.address) }
        when (null == foundDevice) {
            true -> {
                user.devices.add(device.apply { token = "AAA" })
            }

            false -> {
                foundDevice!!.apply { token = "BBB" }
            }
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
        TODO("Not yet implemented")
    }

    override suspend fun getAll(userId: String?): List<Device> {
        TODO("Not yet implemented")
    }

    override suspend fun getByAuth(userId: String?, token: String?): Device {
        TODO("Not yet implemented")
    }
}