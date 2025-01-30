import 'package:grpc/grpc.dart';
import 'package:kiosk/proto/core.pbgrpc.dart';

class Core {
  late final ClientChannel _channel;

  late final client = CoreClient(_channel);

  Core(String uri) {
    final host = uri.split(':')[0];
    final port = int.parse(uri.split(':')[1]);

    _channel = ClientChannel(
      host,
      port: port,
      options: ChannelOptions(
        credentials: ChannelCredentials.insecure(),
      ),
    );
  }

  void emit(String type, {Map<String, String>? data}) {
    client.emit(EmitRequest(type: type, data: data ?? {}));
  }

  void dispose() {
    _channel.shutdown();
  }
}
