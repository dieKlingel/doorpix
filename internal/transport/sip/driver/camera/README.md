# PJSUA2 Camera Driver

The purpose of this driver, is to register a custom video device, which can be integrated with the doorpix camara sessions, to allow sharing a single camera with the rest of the application.

The camera driver creates a virtual camera device which transmitts the stream from the doorpix's internal camera driver comming from gstreamer. The virtual device is named `Internal` and shows up as the driver `github.com/dieklingel/doorpix`.
 