import { bootstrapApplication } from '@angular/platform-browser';
import { provideHttpClient, withFetch } from '@angular/common/http';
import { appConfig } from './app/app.config';
import { AppComponent } from './app/app.component';

// Konfiguriši aplikaciju sa provideHttpClient()
bootstrapApplication(AppComponent, {
  providers: [
    provideHttpClient(withFetch()), // Dodaj provideHttpClient() ovde
    // Dodaj ostale providere ako ih imaš
    ...appConfig.providers // Ako appConfig ima druge providere
  ],
})
  .catch((err) => console.error(err));
