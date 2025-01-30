import 'package:kiosk/config/kiosk_config.dart';
import 'package:yaml/yaml.dart';

class Config {
  final KioskConfig kiosk;

  const Config({
    required this.kiosk,
  });

  factory Config.fromYaml(YamlNode node) {
    if (node is! YamlMap) {
      throw ArgumentError.value(
        node,
        'the node is not a map',
      );
    }
    if (!node.containsKey('kiosk')) {
      throw ArgumentError.value(
        node,
        'the node does not contain the key kiosk',
      );
    }
    if (node['kiosk'] is! YamlMap) {
      throw ArgumentError.value(
        node,
        'the value of the key kiosk is not a map',
      );
    }

    final kiosk = KioskConfig.fromYaml(node['kiosk']);

    return Config(
      kiosk: kiosk,
    );
  }
}
