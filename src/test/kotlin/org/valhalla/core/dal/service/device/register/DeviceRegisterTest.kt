package org.valhalla.core.dal.service.device.register

import kotlinx.coroutines.runBlocking
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.Test
import org.valhalla.core.dal.mock.*
import org.valhalla.core.dal.service.device.DeviceDataAccess
import org.valhalla.core.dal.service.user.UserDataAccess
import org.valhalla.core.dal.tool.BaseDataAccessTest
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
    fun `register and get device (happy path)`(): Unit = runBlocking {

        val user = TestUser.copy()
        user.id = userRepository.register(user)

        val device = TestLinuxDevice.copy()
        device.token = deviceRepository.register(user.id, device)
        val obtainedDevice = deviceRepository.getByAuth(user.id, device.token)

        assertEquals(device, obtainedDevice)
    }


    @Test
    fun `get all devices by user (happy path)`(): Unit = runBlocking {

        val user = TestUser.copy()
        user.id = userRepository.register(user)

        val linuxDevice = TestLinuxDevice.copy()
        val macDevice = TestMacDevice.copy()
        val androidDevice = TestAndroidDevice.copy()
        val iphoneDevice = TestIphoneDevice.copy()

        linuxDevice.token = deviceRepository.register(user.id, linuxDevice)
        macDevice.token = deviceRepository.register(user.id, macDevice)
        androidDevice.token = deviceRepository.register(user.id, androidDevice)
        iphoneDevice.token = deviceRepository.register(user.id, iphoneDevice)

        val obtainedDevices = deviceRepository.getAll(user.id)
        assert(obtainedDevices == listOf(linuxDevice, macDevice, androidDevice, iphoneDevice))

    }
}