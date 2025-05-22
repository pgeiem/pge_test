import os
import sys
import yaml
from datetime import datetime, date

def ensure_isodate(val):
    """Convertit une date/heure en chaîne ISO 8601."""
    if isinstance(val, (datetime, date)):
        return val.isoformat()
    if isinstance(val, str):
        try:
            dt = datetime.fromisoformat(val)
            return dt.isoformat()
        except Exception:
            return val
    return val

def should_keep_desc(desc):
    """Retourne True si le champ desc doit être conservé comme commentaire."""
    if desc is None:
        return False
    if isinstance(desc, str):
        if desc.strip() == "" or desc.strip() == "...":
            return False
    return True

def process_case(old_case):
    """Transforme un cas de test selon le nouveau format et force les dates en ISO."""
    new_case = {
        'name': old_case['name'],
        'now': ensure_isodate(old_case['now']),
        'tests': [{
            'amount': old_case['amount'],
            'end': ensure_isodate(old_case['end']),
        }]
    }
    desc = old_case.get('desc', None)
    if should_keep_desc(desc):
        return new_case, desc
    else:
        return new_case, None

def write_yaml_with_comments(filepath, cases_with_desc):
    with open(filepath, 'w', encoding='utf-8') as f:
        for case, desc in cases_with_desc:
            if desc:
                for line in desc.splitlines():
                    f.write(f"# {line}\n")
            yaml.dump([case], f, allow_unicode=True, sort_keys=False)
            f.write('\n')

def convert_test_files(directory):
    for filename in os.listdir(directory):
        if not filename.endswith('.test'):
            continue
        input_path = os.path.join(directory, filename)
        output_path = os.path.join(directory, filename.replace('.test', '.tests'))
        with open(input_path, 'r', encoding='utf-8') as f:
            old_cases = yaml.safe_load(f)
        cases_with_desc = [process_case(old_case) for old_case in old_cases]
        write_yaml_with_comments(output_path, cases_with_desc)
        print(f"Converti : {filename} -> {os.path.basename(output_path)}")

def main():
    directory = sys.argv[1] if len(sys.argv) > 1 else 'testdata'
    if not os.path.isdir(directory):
        print(f"Le dossier '{directory}' n'existe pas.")
        sys.exit(1)
    convert_test_files(directory)
    print("Conversion terminée.")

if __name__ == "__main__":
    main()
