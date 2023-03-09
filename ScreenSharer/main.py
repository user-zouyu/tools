import os.path

import pyautogui

from pynput import keyboard

import requests
import time

host = "http://localhost:8000/api"
username = "username"

prev_screenshot_time = time.time()


def on_activate_screenshot():
    global prev_screenshot_time
    try:
        if time.time() - prev_screenshot_time > 1.0:
            prev_screenshot_time = time.time()
            image_name = f'img/temp.png'
            pyautogui.screenshot(image_name)
            with open(image_name, "rb") as f:
                requests.post(host + "/upload", params={"username": username}, files={"file": f})
            print(f"发送照片: {image_name}")
    except():
        print("发送失败")


def on_activate_command_prev():
    try:
        requests.get(host + "/command", params={"username": username, "command": -1})
        print(f"发送命令: prev")
    except():
        pass


def on_activate_command_next():
    try:
        requests.get(host + "/command", params={"username": username, "command": 1})
        print(f"发送命令: next")
    except():
        pass


def on_press(key):
    screenshot_hotkey.press(l.canonical(key))
    command_prev_hotkey.press(l.canonical(key))
    command_next_hotkey.press(l.canonical(key))


def on_release(key):
    screenshot_hotkey.release(l.canonical(key))
    command_prev_hotkey.release(l.canonical(key))
    command_next_hotkey.release(l.canonical(key))


if __name__ == "__main__":
    if not os.path.exists("./img"):
        os.mkdir("./img")

    screenshot_hotkey = keyboard.HotKey(
        keyboard.HotKey.parse('<shift>+x'),
        on_activate_screenshot
    )

    command_prev_hotkey = keyboard.HotKey(
        keyboard.HotKey.parse('<shift>+w'),
        on_activate_command_prev
    )

    command_next_hotkey = keyboard.HotKey(
        keyboard.HotKey.parse('<shift>+s'),
        on_activate_command_next)
    with keyboard.Listener(
            on_press=on_press,
            on_release=on_release
    ) as l:
        l.join()
