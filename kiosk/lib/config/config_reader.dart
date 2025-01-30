import 'dart:io';

import 'package:kiosk/config/config.dart';
import 'package:yaml/yaml.dart';

class ConfigReader {
  List<String> _paths = [];

  ConfigReader();

  void addPath(String path) {
    _paths.add(path);
  }

  Future<Config> readConfig() async {
    for (var path in _paths) {
      final file = File(path);
      if (await file.exists()) {
        final content = await file.readAsString();
        final yaml = loadYamlNode(content);
        return Config.fromYaml(yaml);
      }
    }

    throw ArgumentError('no config file found in $_paths');
  }
}
