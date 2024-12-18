package org.valhalla.core.dal.mock

import org.valhalla.core.sdk.model.device.Device
import org.valhalla.core.sdk.model.user.User


val TestUser = User(
    username = "bot.valhalla",
    email = "bot@valhalla.org",
    password = "UBLq<+Z^,g(fr\"S]nu;sCTNX@Wx-pA8m9DtaE{!wkJR5Q}_Pb",
    validated = false,
    validationCode = "087097"
)

val TestLinuxDevice = Device(
    userAgent = "Test/valhalla | Linux | CPU x86",
    address = "10.0.1.1",
    token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
)

val TestMacDevice = Device(
    userAgent = "Test/valhalla | MacOS | CPU AArch64",
    address = "10.0.1.2",
    token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxLjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRP9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf96POk6yJV_adQssw9c"
)

val TestAndroidDevice = Device(
    userAgent = "Test/valhalla | Android | CPU Arm64",
    address = "10.0.1.3",
    token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdUOiOiIxLjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRP9lIiwiaWF0IjoxMYE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf96POk6yJV_adQssw9c"
)

val TestIphoneDevice = Device(
    userAgent = "Test/valhalla | Iphone | CPU Arm64",
    address = "10.0.1.4",
    token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdUOiOiIxLjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRP4JlIiwiaWF0IjoxMYE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf96POk6yPL_adQssw5c"
)