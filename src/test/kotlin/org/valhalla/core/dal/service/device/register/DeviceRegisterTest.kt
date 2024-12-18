package org.valhalla.core.dal.service.device.register

import kotlinx.coroutines.runBlocking
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.Test
import org.valhalla.core.dal.mock.CorrectUser
import org.valhalla.core.dal.service.device.DeviceDataAccess
import org.valhalla.core.dal.service.user.UserDataAccess
import org.valhalla.core.dal.tool.BaseDataAccessTest
import org.valhalla.core.sdk.model.device.Device
import org.valhalla.core.sdk.repository.DeviceRepository
import org.valhalla.core.sdk.repository.UserRepository
import kotlin.test.assertEquals

class DeviceRegisterTest : BaseDataAccessTest() {

    companion object {
        lateinit var userRepository: UserRepository
        lateinit var deviceRepository: DeviceRepository
    }

    @BeforeEach
    fun resetDatabaseImpl() {
        userRepository = UserDataAccess(database!!)
        deviceRepository = DeviceDataAccess(database = database!!, userRepository = userRepository)
    }

    @Test
    fun `register device (happy path)`(): Unit = runBlocking {

        val user = CorrectUser.copy()
        user.id = userRepository.register(user)

        val device = Device(
            userAgent = "Test/valhalla | Linux | CPU x86",
            address = "127.0.0.1",
            token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
        )

        device.token = deviceRepository.register(user.id, device)
        val obtainedDevice = deviceRepository.getByAuth(user.id, device.token)

        assertEquals(device, obtainedDevice)


    }
}