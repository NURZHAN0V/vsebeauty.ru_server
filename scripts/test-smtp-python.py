#!/usr/bin/env python3
"""
Тест SMTP сервера без telnet
Использование: python3 test-smtp-python.py <host> <port> <email-to>
Пример: python3 test-smtp-python.py localhost 2525 test@vsebeauty.ru
"""

import sys
import socket
import time

def test_smtp(host, port, to_email):
    """Тестирует SMTP сервер"""
    print(f"Подключение к SMTP серверу {host}:{port}...")
    
    try:
        # Создаём сокет
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.settimeout(5)
        sock.connect((host, port))
        print("✓ Подключение успешно!")
        
        # Читаем приветствие
        response = sock.recv(1024).decode('utf-8')
        print(f"Приветствие: {response.strip()}")
        
        # EHLO
        sock.send(b"EHLO test\r\n")
        response = sock.recv(1024).decode('utf-8')
        print(f"EHLO: {response.strip()}")
        
        # MAIL FROM
        sock.send(b"MAIL FROM:<test@example.com>\r\n")
        response = sock.recv(1024).decode('utf-8')
        print(f"MAIL FROM: {response.strip()}")
        if not response.startswith(b'250'):
            print("Ошибка в MAIL FROM")
            return False
        
        # RCPT TO
        cmd = f"RCPT TO:<{to_email}>\r\n".encode('utf-8')
        sock.send(cmd)
        response = sock.recv(1024).decode('utf-8')
        print(f"RCPT TO: {response.strip()}")
        if not response.startswith(b'250'):
            print(f"Ошибка в RCPT TO: {response}")
            return False
        
        # DATA
        sock.send(b"DATA\r\n")
        response = sock.recv(1024).decode('utf-8')
        print(f"DATA: {response.strip()}")
        
        # Тело письма
        message = f"""Subject: Test Message

Это тестовое письмо от {time.strftime('%Y-%m-%d %H:%M:%S')}
.
"""
        sock.send(message.encode('utf-8'))
        response = sock.recv(1024).decode('utf-8')
        print(f"Ответ на письмо: {response.strip()}")
        
        # QUIT
        sock.send(b"QUIT\r\n")
        response = sock.recv(1024).decode('utf-8')
        print(f"QUIT: {response.strip()}")
        
        sock.close()
        print("\n✓ Письмо отправлено успешно!")
        print("Проверьте логи сервера и базу данных")
        return True
        
    except socket.timeout:
        print("Ошибка: таймаут подключения")
        return False
    except ConnectionRefusedError:
        print("Ошибка: соединение отклонено. Убедитесь, что SMTP сервер запущен")
        return False
    except Exception as e:
        print(f"Ошибка: {e}")
        return False

if __name__ == "__main__":
    if len(sys.argv) < 4:
        print("Использование: python3 test-smtp-python.py <host> <port> <email-to>")
        print("Пример: python3 test-smtp-python.py localhost 2525 test@vsebeauty.ru")
        sys.exit(1)
    
    host = sys.argv[1]
    port = int(sys.argv[2])
    to_email = sys.argv[3]
    
    test_smtp(host, port, to_email)

