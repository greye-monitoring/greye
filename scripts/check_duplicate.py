import requests

# Lista degli URL da cui prendere i dati
urls = {
    "greye-0.greye-ha": "http://192.168.1.28:32473/api/v1/application/monitor",
    "greye-1.greye-ha": "http://192.168.1.28:30138/api/v1/application/monitor",
    "greye-2.greye-ha": "http://192.168.1.28:30817/api/v1/application/monitor"
}

# Funzione per recuperare i dati JSON da un URL
def get_data_from_url(url):
    try:
        response = requests.get(url)
        response.raise_for_status()  # Controlla se ci sono errori nella risposta
        return response.json()  # Ritorna il contenuto JSON
    except requests.exceptions.RequestException as e:
        print(f"Errore nel recupero dei dati da {url}: {e}")
        return {}

# Funzione per verificare se ci sono chiavi duplicate (nomi delle applicazioni)
def check_duplicate_applications(all_data):
    app_keys = set()  # Set per tenere traccia delle chiavi delle applicazioni (host)
    duplicates = []  # Lista per tenere traccia dei duplicati

    for app_key in all_data:
        if app_key in app_keys:
            duplicates.append(app_key)
        else:
            app_keys.add(app_key)

    return duplicates

# Funzione principale per controllare le applicazioni su tutti gli URL
def main():
    all_data = {}
    total_services = 0  # Variabile per tenere traccia del numero totale di servizi

    # Recupera i dati da tutti gli URL
    for url in urls:
        data = get_data_from_url(urls[url])
        num_services = len(data)

        # Stampa il numero di servizi per ogni URL
        print(f"Numero di servizi da {urls[url]}: {num_services}")

        # Unisce i dati, controllando se la chiave è già presente
        for key, value in data.items():
            if url != value['scheduledApplication'] :
                continue
            if key not in all_data:

                all_data[key] = value
            else:
                print(f"ERRROR Chiave duplicata: {key}")

        total_services += num_services

    # Verifica se ci sono chiavi duplicate
    duplicates = check_duplicate_applications(all_data)

    # Stampa il totale dei servizi
    print(f"\nNumero totale di servizi: {total_services}")

    if duplicates:
        print("Sono state trovate applicazioni duplicate con i seguenti host (chiavi):")
        for duplicate in duplicates:
            print(duplicate)
    else:
        print("Non ci sono applicazioni duplicate.")

if __name__ == "__main__":
    main()
