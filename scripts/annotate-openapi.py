import argparse
import yaml

# Annotate all API Responses with x-go-name attributes


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("spec")

    args = parser.parse_args()
    with open(args.spec, 'r+') as fh:
        spec = yaml.safe_load(fh)

        for response in spec["components"]["responses"]:
            go_name = response[0].upper() + response[1:] + "APIResponse"
            spec["components"]["responses"][response]["x-go-name"] = go_name

        fh.seek(0)
        fh.truncate()

        yaml.safe_dump(spec, fh)


if __name__ == '__main__':
    main()
