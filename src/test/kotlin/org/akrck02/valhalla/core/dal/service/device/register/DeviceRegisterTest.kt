package org.akrck02.valhalla.core.dal.service.device.register

import kotlinx.coroutines.runBlocking
import org.akrck02.valhalla.core.dal.mock.CorrectUser
import org.akrck02.valhalla.core.dal.service.device.DeviceDataAccess
import org.akrck02.valhalla.core.dal.service.user.UserDataAccess
import org.akrck02.valhalla.core.dal.tool.BaseDataAccessTest
import org.akrck02.valhalla.core.sdk.model.device.Device
import org.akrck02.valhalla.core.sdk.repository.DeviceRepository
import org.akrck02.valhalla.core.sdk.repository.UserRepository
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.Test

class DeviceRegisterTest : BaseDataAccessTest() {

    companion object {
        lateinit var userRepository: UserRepository
        lateinit var deviceRepository: DeviceRepository
    }

    @BeforeEach
    fun resetDatabaseImpl() {
        userRepository = UserDataAccess(database!!)
        deviceRepository = DeviceDataAccess(database!!, userRepository)
    }

    @Test
    fun `register device (happy path)`(): Unit = runBlocking {

        val user = CorrectUser.copy()
        user.id = userRepository.register(user)

        val device = Device(
            userAgent = "Test/valhalla | Linux | CPU x86",
            address = "127.0.0.1"
        )

        device.token = deviceRepository.register(user.id, device)
        println(device)
        println()
    }
}