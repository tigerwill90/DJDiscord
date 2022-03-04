FROM openjdk:11

COPY . .

RUN chmod +x entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
CMD ["java", "-Dnogui=true", "-jar", "JMusicBot-0.3.6.jar"]