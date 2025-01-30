import 'package:kiosk/data/label_data.dart';
import 'package:yaml/yaml.dart';

class KioskConfig {
  final String core;
  final Label label;

  const KioskConfig({
    required this.core,
    required this.label,
  });

  factory KioskConfig.fromYaml(YamlNode node) {
    if (node is! YamlMap) {
      throw ArgumentError.value(
        node,
        'the node is not a map',
      );
    }
    if (!node.containsKey('core')) {
      throw ArgumentError.value(
        node,
        'the node does not contain the key core',
      );
    }
    if (node['core'] is! String) {
      throw ArgumentError.value(
        node,
        'the value of the key core is not a string',
      );
    }
    if (!node.containsKey('label')) {
      throw ArgumentError.value(
        node,
        'the node does not contain the key label',
      );
    }
    if (node['label'] is! YamlMap) {
      throw ArgumentError.value(
        node,
        'the value of the key label is not a map',
      );
    }

    final core = node['core'] as String;
    final label = Label.fromYaml(node['label']);

    return KioskConfig(
      core: core,
      label: label,
    );
  }
}
